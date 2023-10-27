// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: installation_user.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const fetchAllocatedComplaintByEmp = `-- name: FetchAllocatedComplaintByEmp :many
select ca.complaint_id as complaint_id,
ca.id as allocation_id,
ci.id as complaint_info_id,
ci.complaint_address as complaint_address,
ci.client_available_date as on_date,
ci.client_available_time_slot as time_slot,
compl.client_id
from complaint_allocations as ca
join complaint_info as ci
on ca.complaint_id = ci.complaint_id
right join complaints as compl
on ca.complaint_id = compl.id
where ca.allocated_to =$1
and ci.complaint_address <> '' and compl.client_id <> ''
order by ca.created_at desc
`

type FetchAllocatedComplaintByEmpRow struct {
	ComplaintID      uuid.UUID      `json:"complaint_id"`
	AllocationID     uuid.UUID      `json:"allocation_id"`
	ComplaintInfoID  uuid.UUID      `json:"complaint_info_id"`
	ComplaintAddress sql.NullString `json:"complaint_address"`
	OnDate           sql.NullTime   `json:"on_date"`
	TimeSlot         sql.NullString `json:"time_slot"`
	ClientID         string         `json:"client_id"`
}

func (q *Queries) FetchAllocatedComplaintByEmp(ctx context.Context, allocatedTo uuid.UUID) ([]FetchAllocatedComplaintByEmpRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchAllocatedComplaintByEmp, allocatedTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FetchAllocatedComplaintByEmpRow{}
	for rows.Next() {
		var i FetchAllocatedComplaintByEmpRow
		if err := rows.Scan(
			&i.ComplaintID,
			&i.AllocationID,
			&i.ComplaintInfoID,
			&i.ComplaintAddress,
			&i.OnDate,
			&i.TimeSlot,
			&i.ClientID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
