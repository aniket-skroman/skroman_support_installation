package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aniket-skroman/skroman_support_installation/apis"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

type InstallationUserRepository interface {
	Init() (context.Context, context.CancelFunc)
	FetchAllocatedComplaintsByEmp(args db.FetchAllocatedComplaintByEmpTodayParams) ([]db.FetchAllocatedComplaintByEmpTodayRow, error)
	FetchAllocatedComplaintsByEmpPending(args db.FetchAllocatedComplaintsByEmpPendingParams) ([]db.FetchAllocatedComplaintsByEmpPendingRow, error)
	CreateComplaintProgress(args db.CreateComplaintProgressParams) (db.ComplaintProgress, error)
	FetchComplaintProgress(args uuid.UUID) ([]db.ComplaintProgress, error)
	DeleteComplaintProgress(progress_id uuid.UUID) (sql.Result, error)
	FetchAllocatedCompletComplaint(allocated_to uuid.UUID) ([]db.FetchAllocatedCompletComplaintRow, error)
}

type installation_user struct {
	db *apis.Store
}

func NewInstallationUserRepository(db *apis.Store) InstallationUserRepository {
	return &installation_user{
		db: db,
	}
}

func (repo *installation_user) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel
}
