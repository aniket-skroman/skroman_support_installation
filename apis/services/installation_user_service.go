package services

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/helper"
	proxycalls "github.com/aniket-skroman/skroman_support_installation/apis/proxy_calls"
	"github.com/aniket-skroman/skroman_support_installation/apis/repositories"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/google/uuid"
)

type InstallationUserService interface {
	FetchAllocatedComplaintByEmp(req dto.FetchAllocatedComplaintRequestDTO) ([]dto.FetchAllocatedComplaintByEmpDTO, error)
	CreateComplaintProgress(req dto.CreateComplaintProgressRequestDTO) (dto.ComplaintProgressDTO, error)
	FetchComplaintProgress(req string) ([]db.ComplaintProgress, error)
	DeleteComplaintProgress(req string) error
	MakeVerificationPendingStatus(complaint_id string) error
	FetchAllocatedCompletComplaint(allocated_to uuid.UUID) ([]dto.FetchAllocatedComplaintByEmpDTO, error)
}

type installation_user struct {
	installation_repo repositories.InstallationUserRepository
	complaint_service ComplaintService
	complaint_repo    repositories.ComplaintRepository
}

func NewInstallationUserService(repo repositories.InstallationUserRepository, complaint_serv ComplaintService,
	complaint_repo repositories.ComplaintRepository,
) InstallationUserService {
	return &installation_user{
		installation_repo: repo,
		complaint_service: complaint_serv,
		complaint_repo:    complaint_repo,
	}
}

func (serv *installation_user) FetchAllocatedComplaintByEmp(req dto.FetchAllocatedComplaintRequestDTO) ([]dto.FetchAllocatedComplaintByEmpDTO, error) {
	var err error
	allocated_obj_id, err := uuid.Parse(req.AllocatedTo)

	if err != nil {
		return nil, err
	}

	var result []dto.FetchAllocatedComplaintByEmpDTO
	current_date := time.Now()
	date := current_date.Format("2006-01-02")
	date_, _ := time.Parse("2006-01-02", date)

	allcation_tag := strings.ToLower(req.AllocationTag)

	if allcation_tag == "today" {

		args := db.FetchAllocatedComplaintByEmpTodayParams{
			AllocatedTo: allocated_obj_id,
			CreatedAt:   date_,
		}
		result, err = serv.fetch_users_today_complaint(args)
	} else if allcation_tag == "complete" {
		result, err = serv.FetchAllocatedCompletComplaint(allocated_obj_id)
	} else {

		args := db.FetchAllocatedComplaintsByEmpPendingParams{
			AllocatedTo: allocated_obj_id,
			CreatedAt:   date_,
		}
		result, err = serv.fetch_users_pending_complaints(args)
	}

	return result, err
}

func (serv *installation_user) fetch_users_today_complaint(args db.FetchAllocatedComplaintByEmpTodayParams) ([]dto.FetchAllocatedComplaintByEmpDTO, error) {
	result, err := serv.installation_repo.FetchAllocatedComplaintsByEmp(args)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}

	complaints := make([]dto.FetchAllocatedComplaintByEmpDTO, len(result))
	wg := sync.WaitGroup{}
	wg.Add(len(result))
	for i, complaint := range result {

		go func(complaint db.FetchAllocatedComplaintByEmpTodayRow, i int) {
			defer wg.Done()
			day_month := fmt.Sprintf("%d %v", complaint.OnDate.Time.Day(), complaint.OnDate.Time.Month())
			complaints[i] = dto.FetchAllocatedComplaintByEmpDTO{
				ComplaintID:     complaint.ComplaintID,
				AllocationID:    complaint.AllocationID,
				ComplaintInfoID: complaint.ComplaintInfoID,
				OnDate:          day_month,
				TimeSlot:        complaint.TimeSlot.String,
				ClientID:        complaint.ClientID,
			}

			client_info, err := serv.complaint_service.Fetch_client_info(complaint.ClientID)

			if err == nil && client_info != nil {
				rv := reflect.ValueOf(client_info)
				complaints[i].ClientInfo = proxycalls.ClientInfoDTO{
					UserName:     fmt.Sprintf("%v", rv.FieldByName("UserName")),
					EmailID:      fmt.Sprintf("%v", rv.FieldByName("Email")),
					MobileNumber: fmt.Sprintf("%v", rv.FieldByName("Contact")),
				}

				complaints[i].ComplaintAddress = fmt.Sprintf("%v", rv.FieldByName("Address"))
			}

			if complaint.Status != "ALLOCATE" && complaint.Status != "INIT" {
				complaints[i].Status = complaint.Status
			}

		}(complaint, i)

	}
	wg.Wait()
	return complaints, nil
}

func (serv *installation_user) fetch_users_pending_complaints(args db.FetchAllocatedComplaintsByEmpPendingParams) ([]dto.FetchAllocatedComplaintByEmpDTO, error) {
	result, err := serv.installation_repo.FetchAllocatedComplaintsByEmpPending(args)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}

	complaints := make([]dto.FetchAllocatedComplaintByEmpDTO, len(result))
	wg := sync.WaitGroup{}
	wg.Add(len(result))
	for i, complaint := range result {

		go func(complaint db.FetchAllocatedComplaintsByEmpPendingRow, i int) {
			defer wg.Done()
			day_month := fmt.Sprintf("%d %v", complaint.OnDate.Time.Day(), complaint.OnDate.Time.Month())
			complaints[i] = dto.FetchAllocatedComplaintByEmpDTO{
				ComplaintID:      complaint.ComplaintID,
				AllocationID:     complaint.AllocationID,
				ComplaintInfoID:  complaint.ComplaintInfoID,
				ComplaintAddress: complaint.ComplaintAddress.String,
				OnDate:           day_month,
				TimeSlot:         complaint.TimeSlot.String,
				ClientID:         complaint.ClientID,
				Status:           complaint.Status,
			}

			client_info, err := serv.complaint_service.Fetch_client_info(complaint.ClientID)
			if err == nil && client_info != nil {
				rv := reflect.ValueOf(client_info)
				complaints[i].ClientInfo = proxycalls.ClientInfoDTO{
					UserName:     fmt.Sprintf("%v", rv.FieldByName("UserName")),
					EmailID:      fmt.Sprintf("%v", rv.FieldByName("Email")),
					MobileNumber: fmt.Sprintf("%v", rv.FieldByName("Contact")),
				}
				complaints[i].ComplaintAddress = fmt.Sprintf("%v", rv.FieldByName("Address"))

			}

		}(complaint, i)

	}
	wg.Wait()
	return complaints, nil
}

func (serv *installation_user) CreateComplaintProgress(req dto.CreateComplaintProgressRequestDTO) (dto.ComplaintProgressDTO, error) {
	// validate id'd
	complaint_id, err := uuid.Parse(req.ComplaintId)

	if err != nil {
		return dto.ComplaintProgressDTO{}, helper.ERR_INVALID_ID
	}

	// get a token id for statement by
	statement_by, err := uuid.Parse(utils.TOKEN_ID)

	if err != nil {
		return dto.ComplaintProgressDTO{}, helper.Err_Something_Wents_Wrong
	}

	args := db.CreateComplaintProgressParams{
		ComplaintID:       complaint_id,
		StatementBy:       statement_by,
		ProgressStatement: req.ProblemStatement,
	}

	result, err := serv.installation_repo.CreateComplaintProgress(args)

	err = helper.Handle_db_err(err)

	if err != nil {
		return dto.ComplaintProgressDTO{}, err
	}

	return dto.ComplaintProgressDTO(result), nil
}

func (serv *installation_user) FetchComplaintProgress(req string) ([]db.ComplaintProgress, error) {
	obj_id, err := uuid.Parse(req)

	if err != nil {
		return nil, helper.ERR_INVALID_ID
	}

	result, err := serv.installation_repo.FetchComplaintProgress(obj_id)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}

	return result, nil
}

func (serv *installation_user) DeleteComplaintProgress(req string) error {
	progress_id, err := uuid.Parse(req)

	if err != nil {
		return helper.ERR_INVALID_ID
	}

	result, err := serv.installation_repo.DeleteComplaintProgress(progress_id)

	if err != nil {
		return err
	}

	a_rows, _ := result.RowsAffected()
	if a_rows == 0 {
		return helper.Err_Delete_Failed
	}

	return nil
}

/* update complaint status when installation user make it complete or resolved, make a status pending-verification */
func (serv *installation_user) MakeVerificationPendingStatus(complaint_id string) error {
	complaint_obj_id, err := uuid.Parse(complaint_id)

	if err != nil {
		return err
	}

	args := db.UpdateComplaintStatusParams{
		ComplaintID: complaint_obj_id,
		Status:      "VERIFICATION-PENDING",
	}

	result, err := serv.complaint_repo.UpdateComplaintStatus(args)

	if err != nil {
		return err
	}

	a_rows, _ := result.RowsAffected()

	if a_rows == 0 {
		return helper.Err_Update_Failed
	}

	return nil
}

func (serv *installation_user) FetchAllocatedCompletComplaint(allocated_to uuid.UUID) ([]dto.FetchAllocatedComplaintByEmpDTO, error) {

	result, err := serv.installation_repo.FetchAllocatedCompletComplaint(allocated_to)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}
	complaints := make([]dto.FetchAllocatedComplaintByEmpDTO, len(result))
	wg := sync.WaitGroup{}
	wg.Add(len(result))
	for i, complaint := range result {

		go func(complaint db.FetchAllocatedCompletComplaintRow, i int) {
			defer wg.Done()
			day_month := fmt.Sprintf("%d %v", complaint.OnDate.Time.Day(), complaint.OnDate.Time.Month())
			complaints[i] = dto.FetchAllocatedComplaintByEmpDTO{
				ComplaintID:     complaint.ComplaintID,
				AllocationID:    complaint.AllocationID,
				ComplaintInfoID: complaint.ComplaintInfoID,
				OnDate:          day_month,
				TimeSlot:        complaint.TimeSlot.String,
				ClientID:        complaint.ClientID,
			}

			client_info, err := serv.complaint_service.Fetch_client_info(complaint.ClientID)

			if err == nil && client_info != nil {
				rv := reflect.ValueOf(client_info)
				complaints[i].ClientInfo = proxycalls.ClientInfoDTO{
					UserName:     fmt.Sprintf("%v", rv.FieldByName("UserName")),
					EmailID:      fmt.Sprintf("%v", rv.FieldByName("Email")),
					MobileNumber: fmt.Sprintf("%v", rv.FieldByName("Contact")),
				}

				complaints[i].ComplaintAddress = fmt.Sprintf("%v", rv.FieldByName("Address"))
			}

			if complaint.Status != "ALLOCATE" && complaint.Status != "INIT" {
				complaints[i].Status = complaint.Status
			}

		}(complaint, i)

	}
	wg.Wait()
	return complaints, nil
}
