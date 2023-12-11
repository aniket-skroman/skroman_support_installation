package repositories

import (
	"database/sql"

	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

func (repo *installation_user) FetchAllocatedComplaintsByEmp(args db.FetchAllocatedComplaintByEmpTodayParams) ([]db.FetchAllocatedComplaintByEmpTodayRow, error) {
	ctx, cancel := repo.Init()
	defer cancel()
	return repo.db.Queries.FetchAllocatedComplaintByEmpToday(ctx, args)
}

func (repo *installation_user) FetchAllocatedComplaintsByEmpPending(args db.FetchAllocatedComplaintsByEmpPendingParams) ([]db.FetchAllocatedComplaintsByEmpPendingRow, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.FetchAllocatedComplaintsByEmpPending(ctx, args)
}

func (repo *installation_user) CreateComplaintProgress(args db.CreateComplaintProgressParams) (db.ComplaintProgress, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CreateComplaintProgress(ctx, args)
}

// this will accept complaintId or statement by
func (repo *installation_user) FetchComplaintProgress(args uuid.UUID) ([]db.ComplaintProgress, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.FetchComplaintProgress(ctx, args)
}

func (repo *installation_user) DeleteComplaintProgress(progress_id uuid.UUID) (sql.Result, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.DeleteComplaintProgressById(ctx, progress_id)
}
