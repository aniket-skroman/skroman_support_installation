// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: installation_user.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createComplaintProgress = `-- name: CreateComplaintProgress :one
insert into complaint_progress (
    complaint_id,
    progress_statement,
    statement_by
) values (
    $1,$2,$3
) returning id, complaint_id, progress_statement, statement_by, created_at, updated_at
`

type CreateComplaintProgressParams struct {
	ComplaintID       uuid.UUID `json:"complaint_id"`
	ProgressStatement string    `json:"progress_statement"`
	StatementBy       uuid.UUID `json:"statement_by"`
}

// add a remarks for complaints , add progress for complaint resolve by installation user
func (q *Queries) CreateComplaintProgress(ctx context.Context, arg CreateComplaintProgressParams) (ComplaintProgress, error) {
	row := q.db.QueryRowContext(ctx, createComplaintProgress, arg.ComplaintID, arg.ProgressStatement, arg.StatementBy)
	var i ComplaintProgress
	err := row.Scan(
		&i.ID,
		&i.ComplaintID,
		&i.ProgressStatement,
		&i.StatementBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteComplaintProgressById = `-- name: DeleteComplaintProgressById :execresult
delete from complaint_progress
where id = $1
`

// delete complaint statement
func (q *Queries) DeleteComplaintProgressById(ctx context.Context, id uuid.UUID) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteComplaintProgressById, id)
}

const fetchAllocatedComplaintByEmpToday = `-- name: FetchAllocatedComplaintByEmpToday :many
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
where ca.allocated_to =$1 and compl.client_id not like '%User_id%'
and ca.created_at < $2 
order by ca.created_at desc
`

type FetchAllocatedComplaintByEmpTodayParams struct {
	AllocatedTo uuid.UUID `json:"allocated_to"`
	CreatedAt   time.Time `json:"created_at"`
}

type FetchAllocatedComplaintByEmpTodayRow struct {
	ComplaintID      uuid.UUID      `json:"complaint_id"`
	AllocationID     uuid.UUID      `json:"allocation_id"`
	ComplaintInfoID  uuid.UUID      `json:"complaint_info_id"`
	ComplaintAddress sql.NullString `json:"complaint_address"`
	OnDate           sql.NullTime   `json:"on_date"`
	TimeSlot         sql.NullString `json:"time_slot"`
	ClientID         string         `json:"client_id"`
}

func (q *Queries) FetchAllocatedComplaintByEmpToday(ctx context.Context, arg FetchAllocatedComplaintByEmpTodayParams) ([]FetchAllocatedComplaintByEmpTodayRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchAllocatedComplaintByEmpToday, arg.AllocatedTo, arg.CreatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FetchAllocatedComplaintByEmpTodayRow{}
	for rows.Next() {
		var i FetchAllocatedComplaintByEmpTodayRow
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

const fetchAllocatedComplaintsByEmpPending = `-- name: FetchAllocatedComplaintsByEmpPending :many
select 
    ca.complaint_id as complaint_id,
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
where ca.allocated_to =$1 and compl.client_id not like '%User_id%'
and ca.created_at < $2 
order by ca.created_at desc
`

type FetchAllocatedComplaintsByEmpPendingParams struct {
	AllocatedTo uuid.UUID `json:"allocated_to"`
	CreatedAt   time.Time `json:"created_at"`
}

type FetchAllocatedComplaintsByEmpPendingRow struct {
	ComplaintID      uuid.UUID      `json:"complaint_id"`
	AllocationID     uuid.UUID      `json:"allocation_id"`
	ComplaintInfoID  uuid.UUID      `json:"complaint_info_id"`
	ComplaintAddress sql.NullString `json:"complaint_address"`
	OnDate           sql.NullTime   `json:"on_date"`
	TimeSlot         sql.NullString `json:"time_slot"`
	ClientID         string         `json:"client_id"`
}

func (q *Queries) FetchAllocatedComplaintsByEmpPending(ctx context.Context, arg FetchAllocatedComplaintsByEmpPendingParams) ([]FetchAllocatedComplaintsByEmpPendingRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchAllocatedComplaintsByEmpPending, arg.AllocatedTo, arg.CreatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FetchAllocatedComplaintsByEmpPendingRow{}
	for rows.Next() {
		var i FetchAllocatedComplaintsByEmpPendingRow
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

const fetchComplaintProgress = `-- name: FetchComplaintProgress :many
select id, complaint_id, progress_statement, statement_by, created_at, updated_at from complaint_progress
where complaint_id = $1 or statement_by = $1
order by created_at desc
`

// fetch all progress by complaint id or empl
func (q *Queries) FetchComplaintProgress(ctx context.Context, complaintID uuid.UUID) ([]ComplaintProgress, error) {
	rows, err := q.db.QueryContext(ctx, fetchComplaintProgress, complaintID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ComplaintProgress{}
	for rows.Next() {
		var i ComplaintProgress
		if err := rows.Scan(
			&i.ID,
			&i.ComplaintID,
			&i.ProgressStatement,
			&i.StatementBy,
			&i.CreatedAt,
			&i.UpdatedAt,
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
