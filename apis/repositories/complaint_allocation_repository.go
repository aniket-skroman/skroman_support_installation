package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aniket-skroman/skroman_support_installation/apis"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

type ComplaintAllocationRepository interface {
	Init() (context.Context, context.CancelFunc)
	CreateComplaintAllocation(db.CreateComplaintAllocationParams) error
	UpdateComplaintAllocation(db.UpdateComplaintAllocationParams) error
	FetchAllocationByComplaintId(uuid.UUID) (db.ComplaintAllocations, error)
	CheckComplaintStatusBeforeUpdate(uuid.UUID) (sql.Result, error)
	CheckDuplicateComplaintAllocation(uuid.UUID) (sql.Result, error)
}

type allocation_repository struct {
	db *apis.Store
}

func NewComplaintAllocationRepository(db *apis.Store) ComplaintAllocationRepository {
	return &allocation_repository{
		db: db,
	}
}

func (repo *allocation_repository) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel
}
