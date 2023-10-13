package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/repositories"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

type ComplaintService interface {
	CreateComplaint(dto.CreateComplaintRequestDTO) (interface{}, error)
	FetchAllComplaints(dto.PaginationRequestParams) ([]dto.ComplaintInfoDTO, error)
}

type complaint_service struct {
	complaint_repo repositories.ComplaintRepository
}

func NewComplaintService(repo repositories.ComplaintRepository) ComplaintService {
	return &complaint_service{
		complaint_repo: repo,
	}
}

func (ser *complaint_service) CreateComplaint(req dto.CreateComplaintRequestDTO) (interface{}, error) {
	//init a complaint first
	comp_args := db.CreateComplaintParams{
		ClientID:  req.ClientID,
		CreatedBy: uuid.New(),
	}
	complaint, err := ser.complaint_repo.CreateComplaint(comp_args)

	if err != nil {
		return nil, err
	}

	// create a complaint info
	c_time, err := time.Parse("2006-01-02 15:04:05", req.ClientAvailable)
	if err != nil {
		return nil, err
	}
	info_args := db.CreateComplaintInfoParams{
		ComplaintID: complaint.ID,
		DeviceID:    req.DeviceID,
		DeviceType: sql.NullString{
			String: req.DeviceModel, Valid: true,
		},
		DeviceModel:      sql.NullString{String: req.DeviceModel, Valid: true},
		ProblemStatement: req.ProblemStatement,
		ProblemCategory:  sql.NullString{String: req.ProblemCategory, Valid: true},
		ClientAvailable:  c_time,
		Status:           "INIT",
	}

	complaint_info, err := ser.complaint_repo.CreateComplaintInfo(info_args)

	complaint_info_dto := dto.ComplaintInfoDTO{
		ID:               complaint_info.ID,
		ComplaintID:      complaint_info.ComplaintID,
		DeviceID:         complaint_info.DeviceID,
		ProblemStatement: complaint_info.ProblemStatement,
		ProblemCategory:  complaint_info.ProblemCategory.String,
		ClientAvailable:  complaint_info.ClientAvailable,
		Status:           complaint_info.Status,
		CreatedAt:        complaint_info.CreatedAt,
		UpdatedAt:        complaint_info.UpdatedAt,
		DeviceType:       complaint_info.DeviceType.String,
		DeviceModel:      complaint_info.DeviceModel.String,
	}

	return complaint_info_dto, err
}

func (ser *complaint_service) FetchAllComplaints(req dto.PaginationRequestParams) ([]dto.ComplaintInfoDTO, error) {
	args := db.FetchAllComplaintsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	complaints, err := ser.complaint_repo.FetchAllComplaints(args)

	if err != nil {
		return nil, err
	}

	if len(complaints) == 0 {
		return nil, errors.New("complaint not found")
	}

	complaint_info := new(dto.ComplaintInfoDTO).SetComplaintInfoData(complaints...)

	if _, ok := complaint_info.([]dto.ComplaintInfoDTO); ok {
		return complaint_info.([]dto.ComplaintInfoDTO), nil
	}

	s_complaint := complaint_info.(dto.ComplaintInfoDTO)

	return []dto.ComplaintInfoDTO{s_complaint}, nil

}
