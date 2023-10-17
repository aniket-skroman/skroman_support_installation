package dto

import (
	"time"

	proxycalls "github.com/aniket-skroman/skroman_support_installation/apis/proxy_calls"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

type ClientTimeSlots struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}

type CreateComplaintRequestDTO struct {
	ClientID            string          `json:"client_id" binding:"required"`
	DeviceID            string          `json:"device_id" binding:"required"`
	DeviceType          string          `json:"device_type" binding:"required"`
	DeviceModel         string          `json:"device_model" binding:"required"`
	ProblemStatement    string          `json:"problem_statement" binding:"required"`
	ProblemCategory     string          `json:"problem_category" binding:"required"`
	ClientAvailable     string          `json:"client_available" binding:"required"`
	ClientAvailableDate string          `json:"client_available_date" binding:"required"`
	ClientTimeSlots     ClientTimeSlots `json:"available_time_slots" binding:"required"`
}

type PaginationRequestParams struct {
	PageID   int32 `uri:"page_id"`
	PageSize int32 `uri:"page_size"`
}

type ImageRequestDTO struct {
	ImagePath string `uri:"image_path"`
}

type ComplaintDeviceImagesDTO struct {
	DeviceImage string    `json:"device_image"`
	CreatedAt   time.Time `json:"uploaded_at"`
}

type ComplaintInfoByComplaintDTO struct {
	ClientInfo    proxycalls.ClientInfoDTO `json:"client_info"`
	ComplaintInfo ComplaintFullDetailsDTO  `json:"complaint_info"`
}

type ComplaintFullDetailsDTO struct {
	CreatedBy         string                     `json:"created_by"`
	Client            string                     `json:"client"`
	DeviceID          string                     `json:"device_id"`
	ProblemStatement  string                     `json:"problem_statement"`
	ProblemCategory   string                     `json:"problem_category"`
	ClientAvailable   time.Time                  `json:"client_available"`
	ComplaintStatus   string                     `json:"complaint_status"`
	DeviceModel       string                     `json:"device_model"`
	DeviceType        string                     `json:"device_type"`
	ComplaintRaisedAt time.Time                  `json:"complaint_raised_at"`
	LastModifiedAt    time.Time                  `json:"last_modified_at"`
	DeviceInfo        []ComplaintDeviceImagesDTO `json:"device_images"`
}

type ComplaintInfoDTO struct {
	ID               uuid.UUID `json:"id"`
	ComplaintID      uuid.UUID `json:"complaint_id"`
	DeviceID         string    `json:"device_id"`
	ProblemStatement string    `json:"problem_statement"`
	ProblemCategory  string    `json:"problem_category"`
	ClientAvailable  time.Time `json:"client_available"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeviceType       string    `json:"device_type"`
	DeviceModel      string    `json:"device_model"`
}

func (complaint *ComplaintInfoDTO) SetComplaintInfoData(module_data ...db.ComplaintInfo) interface{} {
	if len(module_data) == 1 {
		return ComplaintInfoDTO{
			ID:               module_data[0].ID,
			ComplaintID:      module_data[0].ComplaintID,
			DeviceID:         module_data[0].DeviceID,
			ProblemStatement: module_data[0].ProblemStatement,
			ProblemCategory:  module_data[0].ProblemCategory.String,
			ClientAvailable:  module_data[0].ClientAvailable,
			Status:           module_data[0].Status,
			CreatedAt:        module_data[0].CreatedAt,
			UpdatedAt:        module_data[0].UpdatedAt,
			DeviceType:       module_data[0].DeviceType.String,
			DeviceModel:      module_data[0].DeviceModel.String,
		}
	}

	complaints := make([]ComplaintInfoDTO, len(module_data))

	for i := range module_data {
		complaints[i] = ComplaintInfoDTO{
			ID:               module_data[i].ID,
			ComplaintID:      module_data[i].ComplaintID,
			DeviceID:         module_data[i].DeviceID,
			ProblemStatement: module_data[i].ProblemStatement,
			ProblemCategory:  module_data[i].ProblemCategory.String,
			ClientAvailable:  module_data[i].ClientAvailable,
			Status:           module_data[i].Status,
			CreatedAt:        module_data[i].CreatedAt,
			UpdatedAt:        module_data[i].UpdatedAt,
			DeviceType:       module_data[i].DeviceType.String,
			DeviceModel:      module_data[i].DeviceModel.String,
		}
	}

	return complaints
}
