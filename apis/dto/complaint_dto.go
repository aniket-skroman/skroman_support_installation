package dto

import (
	"strings"
	"time"

	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

type ClientTimeSlots struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}

type CreateComplaintRequestDTO struct {
	ClientID            string          `json:"client_id" binding:"required"`
	DeviceID            string          `json:"device_id"`
	DeviceType          string          `json:"device_type" binding:"required"`
	DeviceModel         string          `json:"device_model" binding:"required"`
	ProblemStatement    string          `json:"problem_statement" binding:"required"`
	ProblemCategory     string          `json:"problem_category" binding:"required"`
	ClientAvailableDate string          `json:"client_available_date" binding:"required"`
	ClientTimeSlots     ClientTimeSlots `json:"available_time_slots" binding:"required"`
}

type UpdateComplaintRequestDTO struct {
	ComplaintInfoId     string          `json:"complaint_info_id" binding:"required"`
	DeviceID            string          `json:"device_id"`
	DeviceType          string          `json:"device_type" binding:"required"`
	DeviceModel         string          `json:"device_model" binding:"required"`
	ProblemStatement    string          `json:"problem_statement" binding:"required"`
	ProblemCategory     string          `json:"problem_category" binding:"required"`
	ClientAvailableDate string          `json:"client_available_date" binding:"required"`
	ClientTimeSlots     ClientTimeSlots `json:"available_time_slots" binding:"required"`
}

type PaginationRequestParams struct {
	PageID   int32  `uri:"page_id" binding:"required"`
	PageSize int32  `uri:"page_size" binding:"required"`
	TagKey   string `uri:"tag_key" binding:"required,oneof=INIT ALLOCATE COMPLETE TOTAL"`
}

type ImageRequestDTO struct {
	ImagePath string `uri:"image_path"`
}

type VideoRequestDTO struct {
	Directory string `uri:"directory"`
	FilePath  string `uri:"image_path"`
}

type ComplaintDeviceImagesDTO struct {
	ID        uuid.UUID `json:"id"`
	File      string    `json:"file"`
	CreatedAt time.Time `json:"uploaded_at"`
	FileType  string    `json:"file_type"`
}

type AllocatedEmpDetailsDTO struct {
	ID             string      `json:"id,omitempty"`
	AllocateUserID string      `json:"allocate_user_id,omitempty"`
	FullName       string      `json:"full_name,omitempty"`
	Email          string      `json:"email,omitempty"`
	Contact        string      `json:"contact,omitempty"`
	UserType       string      `json:"user_type,omitempty"`
	CreatedAt      interface{} `json:"allocate_date,omitempty"`
	UpdatedAt      interface{} `json:"allocate_modify_date,omitempty"`
}

type MonthWiseCounts struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type AllComplaintsCount struct {
	AllComplaints       int64             `json:"all_complaints"`
	PendingComplaints   int64             `json:"pending_complaints"`
	CompletedComplaints int64             `json:"completed_complaints"`
	AllocatedComplaints int64             `json:"allocated_complaints"`
	EmpCount            int64             `json:"client_count"`
	MonthWiseCounts     []MonthWiseCounts `json:"month_data"`
}

type ComplaintInfoByComplaintDTO struct {
	ClientInfo    interface{}             `json:"client_info"`
	ComplaintInfo ComplaintFullDetailsDTO `json:"complaint_info"`
}

type ComplaintFullDetailsDTO struct {
	Id                     uuid.UUID                  `json:"complaint_info_id"`
	ComplaintId            uuid.UUID                  `json:"complaint_id"`
	CreatedBy              string                     `json:"created_by"`
	Client                 string                     `json:"client"`
	DeviceID               string                     `json:"device_id"`
	ProblemStatement       string                     `json:"problem_statement"`
	ProblemCategory        string                     `json:"problem_category"`
	ClientAvailableDate    string                     `json:"client_available_date" binding:"required"`
	ClientTimeSlots        string                     `json:"available_time_slots" binding:"required"`
	ComplaintStatus        string                     `json:"complaint_status"`
	ComplaintAddress       string                     `json:"complaint_address"`
	DeviceModel            string                     `json:"device_model"`
	DeviceType             string                     `json:"device_type"`
	ComplaintRaisedAt      time.Time                  `json:"complaint_raised_at"`
	LastModifiedAt         time.Time                  `json:"last_modified_at"`
	DeviceImages           []ComplaintDeviceImagesDTO `json:"device_images"`
	DeviceVideos           []ComplaintDeviceImagesDTO `json:"device_videos"`
	AllocatedEmpDetailsDTO AllocatedEmpDetailsDTO     `json:"allocation_info"`
}

type ComplaintByClientDTO struct {
	ClientInfo    interface{}        `json:"client_info"`
	ComplaintInfo []ComplaintInfoDTO `json:"complaint_info"`
}

type ComplaintInfoDTO struct {
	ID                  uuid.UUID `json:"id"`
	ComplaintID         uuid.UUID `json:"complaint_id"`
	DeviceID            string    `json:"device_id"`
	ProblemStatement    string    `json:"problem_statement"`
	ProblemCategory     string    `json:"problem_category"`
	ClientAvailableDate string    `json:"client_available_date" binding:"required"`
	ClientTimeSlots     string    `json:"available_time_slots" binding:"required"`
	ComplaintAddress    string    `json:"complaint_address"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeviceType          string    `json:"device_type"`
	DeviceModel         string    `json:"device_model"`
}

func (complaint *ComplaintInfoDTO) SetComplaintInfoData(module_data ...db.ComplaintInfo) interface{} {

	if len(module_data) == 1 {
		a_date := strings.ReplaceAll(module_data[0].ClientAvailableDate.Time.String(), "00:00:00 +0000 UTC", "")

		return ComplaintInfoDTO{
			ID:                  module_data[0].ID,
			ComplaintID:         module_data[0].ComplaintID,
			DeviceID:            module_data[0].DeviceID,
			ProblemStatement:    module_data[0].ProblemStatement,
			ProblemCategory:     module_data[0].ProblemCategory.String,
			ClientAvailableDate: a_date,
			ClientTimeSlots:     module_data[0].ClientAvailableTimeSlot.String,
			ComplaintAddress:    module_data[0].ComplaintAddress.String,
			Status:              module_data[0].Status,
			CreatedAt:           module_data[0].CreatedAt,
			UpdatedAt:           module_data[0].UpdatedAt,
			DeviceType:          module_data[0].DeviceType.String,
			DeviceModel:         module_data[0].DeviceModel.String,
		}
	}

	complaints := make([]ComplaintInfoDTO, len(module_data))

	for i := range module_data {
		a_date := strings.ReplaceAll(module_data[i].ClientAvailableDate.Time.String(), "00:00:00 +0000 UTC", "")

		complaints[i] = ComplaintInfoDTO{
			ID:                  module_data[i].ID,
			ComplaintID:         module_data[i].ComplaintID,
			DeviceID:            module_data[i].DeviceID,
			ProblemStatement:    module_data[i].ProblemStatement,
			ProblemCategory:     module_data[i].ProblemCategory.String,
			ClientAvailableDate: a_date,
			ClientTimeSlots:     module_data[i].ClientAvailableTimeSlot.String,
			Status:              module_data[i].Status,
			CreatedAt:           module_data[i].CreatedAt,
			UpdatedAt:           module_data[i].UpdatedAt,
			DeviceType:          module_data[i].DeviceType.String,
			DeviceModel:         module_data[i].DeviceModel.String,
			ComplaintAddress:    module_data[i].ComplaintAddress.String,
		}
	}

	return complaints
}

// ----------------------------------------- CLIENT REGISTRATION  ------------------------------------------------ //
type ClientRegistration struct {
	UserName    string `json:"userName" binding:"required"`
	Email       string `json:"emailId" binding:"required,email"`
	Contact     string `json:"mobileNumber" binding:"required,min=10"`
	AddressLine string `json:"address1" binding:"required"`
	City        string `json:"city" binding:"required"`
	State       string `json:"state" binding:"required"`
	Pincode     string `json:"pinCode" binding:"required"`
}

type FetchComplaintsByClientRequestDTO struct {
	ClientId string `uri:"client_id"`
	PageID   int32  `uri:"page_id" binding:"required"`
	PageSize int32  `uri:"page_size" binding:"required"`
}

// ------------------------------------------ COMPLAINT ALLOCATIONS ----------------------------------------------- //
type CreateAllocationRequestDTO struct {
	ComplaintId string `json:"complaint_id" binding:"required"`
	AllocateBy  string `json:"allocate_by" `
	AllocateTo  string `json:"allocate_to" binding:"required"`
}

type UpdateAllocateComplaintRequestDTO struct {
	Id         string `form:"id" binding:"required"`
	AllocateTo string `form:"allocate_to" binding:"required"`
	AllocateBy string `form:"allocate_by" `
}
