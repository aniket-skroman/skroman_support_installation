package services

import (
	"database/sql"
	"sync"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/repositories"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

type InstallationUserService interface {
	FetchAllocatedComplaintByEmp(allocated_id string) ([]dto.FetchAllocatedComplaintByEmpDTO, error)
}

type installation_user struct {
	installation_repo repositories.InstallationUserRepository
	complaint_service ComplaintService
}

func NewInstallationUserService(repo repositories.InstallationUserRepository, complaint_serv ComplaintService) InstallationUserService {
	return &installation_user{
		installation_repo: repo,
		complaint_service: complaint_serv,
	}
}

func (serv *installation_user) FetchAllocatedComplaintByEmp(allocated_id string) ([]dto.FetchAllocatedComplaintByEmpDTO, error) {
	allocated_obj_id, err := uuid.Parse(allocated_id)

	if err != nil {
		return nil, err
	}

	result, err := serv.installation_repo.FetchAllocatedComplaintsByEmp(allocated_obj_id)

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

		go func(complaint db.FetchAllocatedComplaintByEmpRow, i int) {
			defer wg.Done()
			complaints[i] = dto.FetchAllocatedComplaintByEmpDTO{
				ComplaintID:      complaint.ComplaintID,
				AllocationID:     complaint.AllocationID,
				ComplaintInfoID:  complaint.ComplaintInfoID,
				ComplaintAddress: complaint.ComplaintAddress.String,
				OnDate:           complaint.OnDate.Time,
				TimeSlot:         complaint.TimeSlot.String,
				ClientID:         complaint.ClientID,
			}

			// client_info, err := serv.complaint_service.Fetch_client_info(complaint.ClientID)
			// if err == nil && client_info != nil {
			// 	fmt.Println("Client Info Result : ", client_info)
			// }
			// if err == nil {
			// 	// complaints[i].ClientInfo = proxycalls.ClientInfoDTO{
			// 	// 	UserName:     client_info["user_name"].(string),
			// 	// 	EmailID:      client_info.Result.EmailID,
			// 	// 	MobileNumber: client_info.Result.MobileNumber,
			// 	// }
			// }

		}(complaint, i)

	}
	wg.Wait()
	return complaints, nil
}
