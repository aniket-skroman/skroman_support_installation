package dto

import (
	"time"

	proxycalls "github.com/aniket-skroman/skroman_support_installation/apis/proxy_calls"
	"github.com/google/uuid"
)

type FetchAllocatedComplaintByEmpDTO struct {
	ClientInfo       proxycalls.ClientInfoDTO `json:"client_info"`
	ComplaintID      uuid.UUID                `json:"complaint_id"`
	AllocationID     uuid.UUID                `json:"allocation_id"`
	ComplaintInfoID  uuid.UUID                `json:"complaint_info_id"`
	ComplaintAddress string                   `json:"complaint_address"`
	OnDate           string                   `json:"client_available_date"`
	TimeSlot         string                   `json:"time_slot"`
	ClientID         string                   `json:"client_id"`
}

type FetchAllocatedComplaintRequestDTO struct {
	AllocatedTo   string `form:"allocated_to" binding:"required"`
	AllocationTag string `form:"allocation_tag" binding:"required,oneof=Today Pending pending today"`
}

type CreateComplaintProgressRequestDTO struct {
	ComplaintId      string `json:"complaint_id" binding:"required"`
	ProblemStatement string `json:"problem_statement" binding:"required"`
	StatementBy      uuid.UUID
}

type ComplaintProgressDTO struct {
	ID                uuid.UUID `json:"id"`
	ComplaintID       uuid.UUID `json:"complaint_id"`
	ProgressStatement string    `json:"progress_statement"`
	StatementBy       uuid.UUID `json:"statement_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
