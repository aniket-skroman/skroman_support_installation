package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/helper"
	proxycalls "github.com/aniket-skroman/skroman_support_installation/apis/proxy_calls"
	"github.com/aniket-skroman/skroman_support_installation/apis/repositories"
	"github.com/aniket-skroman/skroman_support_installation/notifications"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

type ComplaintAllocationService interface {
	AllocateComplaint(req dto.CreateAllocationRequestDTO) error
	UpdateComplaintAllocation(req dto.UpdateAllocateComplaintRequestDTO) error
	FetchAllocationByComplaintId(complaint_id string) (db.ComplaintAllocations, error)
}

type allocation_service struct {
	allocation_repo repositories.ComplaintAllocationRepository
	jwt_service     JWTService
}

func NewComplaintAllocationService(repo repositories.ComplaintAllocationRepository, jwt_serv JWTService) ComplaintAllocationService {
	return &allocation_service{
		allocation_repo: repo,
		jwt_service:     jwt_serv,
	}
}

func (ser *allocation_service) AllocateComplaint(req dto.CreateAllocationRequestDTO) error {
	var err error
	complaint_obj, err := helper.ValidateUUID(req.ComplaintId)
	if err != nil {
		return helper.ERR_INVALID_ID
	}

	allocate_by, err := helper.ValidateUUID(req.AllocateBy)
	if err != nil {
		return helper.ERR_INVALID_ID
	}

	allocate_to, err := helper.ValidateUUID(req.AllocateTo)
	if err != nil {
		return helper.ERR_INVALID_ID
	}

	do_allocation := make(chan bool)
	is_allocation_done := make(chan bool)
	err_chan := make(chan error)
	wg := sync.WaitGroup{}
	var result sql.Result

	wg.Add(3)

	// check this complaint should not allocated before
	go func() {
		defer wg.Done()
		// check this complaint should not allocated before
		result, err = ser.allocation_repo.CheckDuplicateComplaintAllocation(complaint_obj)
		err_chan <- err
		affected_rows, err := result.RowsAffected()
		err_chan <- err
		if affected_rows != 0 {
			err_chan <- errors.New("this complaint is already allocated to someone")
			return
		}
		do_allocation <- true
	}()

	// fetch the fcm token to notify that user
	go func() {
		for {
			if _, ok := <-is_allocation_done; !ok {
				break
			} else {
				ser.notify_user(allocate_to, &wg)
				close(is_allocation_done)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			if _, ok := <-do_allocation; !ok {
				break
			} else {
				args := db.CreateComplaintAllocationParams{
					ComplaintID: complaint_obj,
					AllocatedBy: allocate_by,
					AllocatedTo: allocate_to,
				}

				err = ser.allocation_repo.CreateComplaintAllocation(args)
				err = helper.Handle_db_err(err)
				if err != nil {
					err_chan <- err
					return
				}
				is_allocation_done <- true
				close(do_allocation)
			}

		}
	}()

	go func() {
		wg.Wait()
		close(err_chan)
	}()

	for data_err := range err_chan {
		if data_err != nil {
			close(do_allocation)
			close(is_allocation_done)
			return data_err
		}
	}

	return nil
}

func (ser *allocation_service) notify_user(user_id uuid.UUID, p_wg *sync.WaitGroup) {
	defer p_wg.Done()
	tokens, _ := ser.fetch_fcm_tokens(user_id)
	n := notifications.Notification{}
	app, _, _ := n.SetupFirebase()

	n.MsgTitle = "Skroman Complaint Allocation"
	n.MsgBody = "New Complaint has been allocated, please check"

	wg := sync.WaitGroup{}

	for i := range tokens {
		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			n.RegistrationToken = token
			n.SendToToken(app)
		}(tokens[i])
	}

	wg.Wait()

}

func (ser *allocation_service) UpdateComplaintAllocation(req dto.UpdateAllocateComplaintRequestDTO) error {
	id, err := helper.ValidateUUID(req.Id)

	if err != nil {
		return helper.ERR_INVALID_ID
	}
	allocate_to, err := helper.ValidateUUID(req.AllocateTo)
	if err != nil {
		return helper.ERR_INVALID_ID
	}

	allocate_by, err := helper.ValidateUUID(req.AllocateBy)
	if err != nil {
		return helper.ERR_INVALID_ID
	}

	// check the complaint status, it should be init/allocate should not be complete
	result, err := ser.allocation_repo.CheckComplaintStatusBeforeUpdate(id)

	if err != nil {
		return err
	}

	affected_rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected_rows == 0 {
		return errors.New("completed complaint will not be updated")
	}

	args := db.UpdateComplaintAllocationParams{
		ID:          id,
		AllocatedTo: allocate_to,
		AllocatedBy: allocate_by,
	}

	err = ser.allocation_repo.UpdateComplaintAllocation(args)

	if err != nil {
		return err
	}

	err = helper.Handle_db_err(err)
	return err
}

func (ser *allocation_service) FetchAllocationByComplaintId(complaint_id string) (db.ComplaintAllocations, error) {
	complaint_obj_id, err := uuid.Parse(complaint_id)

	if err != nil {
		return db.ComplaintAllocations{}, nil
	}

	return ser.allocation_repo.FetchAllocationByComplaintId(complaint_obj_id)
}

func (ser *allocation_service) fetch_fcm_tokens(user_id uuid.UUID) ([]string, error) {
	token := ser.jwt_service.GenerateTempToken(user_id.String(), "EMP", "INSTALLATION")
	end_point := "fcm-data/" + user_id.String()
	proxycall := proxycalls.NewAPIRequest(end_point, http.MethodGet, false, nil, user_id, map[string]string{"Authorization": token})
	response, err := proxycall.MakeApiRequest()

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {

		response_data := struct {
			Error    string `json:"error"`
			Message  string `json:"message"`
			Status   bool   `json:"status"`
			UserData []struct {
				ID        string    `json:"id"`
				UserID    string    `json:"user_id"`
				FcmToken  string    `json:"fcm_token"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
			} `json:"user_data"`
		}{}

		err := json.NewDecoder(response.Body).Decode(&response_data)
		if err != nil {
			return nil, err
		}

		tokens := make([]string, len(response_data.UserData))

		for i, user_data := range response_data.UserData {
			tokens[i] = user_data.FcmToken
			//fmt.Println("Data has been appended..\n", i, user_data.FcmToken)
		}

		return tokens, nil

	} else {
		fmt.Println("Response has not fetched success")
	}

	return nil, nil
}
