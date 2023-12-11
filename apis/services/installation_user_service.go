package services

import (
	"database/sql"
	"fmt"
	"reflect"
	"sync"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	proxycalls "github.com/aniket-skroman/skroman_support_installation/apis/proxy_calls"
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
			day_month := fmt.Sprintf("%d %v", complaint.OnDate.Time.Day(), complaint.OnDate.Time.Month())
			complaints[i] = dto.FetchAllocatedComplaintByEmpDTO{
				ComplaintID:      complaint.ComplaintID,
				AllocationID:     complaint.AllocationID,
				ComplaintInfoID:  complaint.ComplaintInfoID,
				ComplaintAddress: complaint.ComplaintAddress.String,
				OnDate:           day_month,
				TimeSlot:         complaint.TimeSlot.String,
				ClientID:         complaint.ClientID,
			}

			client_info, err := serv.complaint_service.Fetch_client_info(complaint.ClientID)
			if err == nil && client_info != nil {
				rv := reflect.ValueOf(client_info)
				complaints[i].ClientInfo = proxycalls.ClientInfoDTO{
					UserName:     fmt.Sprintf("%v", rv.FieldByName("UserName")),
					EmailID:      fmt.Sprintf("%v", rv.FieldByName("Email")),
					MobileNumber: fmt.Sprintf("%v", rv.FieldByName("Contact")),
				}

			}

		}(complaint, i)

	}
	wg.Wait()
	return complaints, nil
}
