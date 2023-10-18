package repositories

import (
	"context"
	"time"

	"github.com/aniket-skroman/skroman_support_installation/apis"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
)

type ComplaintAllocationRepository interface {
	Init() (context.Context, context.CancelFunc)
	CreateComplaintAllocation(args db.CreateComplaintAllocationParams) error
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