package repositories

import (
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
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
