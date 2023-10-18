// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: complaint_allocation.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createComplaintAllocation = `-- name: CreateComplaintAllocation :one
insert into complaint_allocations (
    complaint_id,
    allocated_to,
    allocated_by
) values (
    $1,$2,$3
) returning id, complaint_id, allocated_to, allocated_by, created_at, updated_at
`

type CreateComplaintAllocationParams struct {
	ComplaintID uuid.UUID `json:"complaint_id"`
	AllocatedTo uuid.UUID `json:"allocated_to"`
	AllocatedBy uuid.UUID `json:"allocated_by"`
}

func (q *Queries) CreateComplaintAllocation(ctx context.Context, arg CreateComplaintAllocationParams) (ComplaintAllocations, error) {
	row := q.db.QueryRowContext(ctx, createComplaintAllocation, arg.ComplaintID, arg.AllocatedTo, arg.AllocatedBy)
	var i ComplaintAllocations
	err := row.Scan(
		&i.ID,
		&i.ComplaintID,
		&i.AllocatedTo,
		&i.AllocatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const fetchComplaintAllocationByComplaint = `-- name: FetchComplaintAllocationByComplaint :one
select id, complaint_id, allocated_to, allocated_by, created_at, updated_at from complaint_allocations
where complaint_id = $1
`

func (q *Queries) FetchComplaintAllocationByComplaint(ctx context.Context, complaintID uuid.UUID) (ComplaintAllocations, error) {
	row := q.db.QueryRowContext(ctx, fetchComplaintAllocationByComplaint, complaintID)
	var i ComplaintAllocations
	err := row.Scan(
		&i.ID,
		&i.ComplaintID,
		&i.AllocatedTo,
		&i.AllocatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateComplaintAllocation = `-- name: UpdateComplaintAllocation :one
update complaint_allocations
set allocated_to=$2,
allocated_by=$3,
updated_at = CURRENT_TIMESTAMP
where id=$1
returning id, complaint_id, allocated_to, allocated_by, created_at, updated_at
`

type UpdateComplaintAllocationParams struct {
	ID          uuid.UUID `json:"id"`
	AllocatedTo uuid.UUID `json:"allocated_to"`
	AllocatedBy uuid.UUID `json:"allocated_by"`
}

func (q *Queries) UpdateComplaintAllocation(ctx context.Context, arg UpdateComplaintAllocationParams) (ComplaintAllocations, error) {
	row := q.db.QueryRowContext(ctx, updateComplaintAllocation, arg.ID, arg.AllocatedTo, arg.AllocatedBy)
	var i ComplaintAllocations
	err := row.Scan(
		&i.ID,
		&i.ComplaintID,
		&i.AllocatedTo,
		&i.AllocatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}