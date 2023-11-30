package repositories

import (
	"context"
	"time"

	"github.com/aniket-skroman/skroman_support_installation/apis"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

type InstallationUserRepository interface {
	Init() (context.Context, context.CancelFunc)
	FetchAllocatedComplaintsByEmp(allocated_id uuid.UUID) ([]db.FetchAllocatedComplaintByEmpRow, error)
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