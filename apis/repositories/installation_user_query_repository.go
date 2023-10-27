package repositories

import (
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

func (repo *installation_user) FetchAllocatedComplaintsByEmp(allocated_id uuid.UUID) ([]db.FetchAllocatedComplaintByEmpRow, error) {
	ctx, cancel := repo.Init()
	defer cancel()
	return repo.db.Queries.FetchAllocatedComplaintByEmp(ctx, allocated_id)
}
