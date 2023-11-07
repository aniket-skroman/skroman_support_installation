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
	OnDate           time.Time                `json:"client_available_date"`
	TimeSlot         string                   `json:"time_slot"`
	ClientID         string                   `json:"client_id"`
}
