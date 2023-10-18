package repositories

import (
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
)

func (repo *allocation_repository) CreateComplaintAllocation(args db.CreateComplaintAllocationParams) error {
	ctx, cancel := repo.Init()
	defer cancel()

	_, err := repo.db.Queries.CreateComplaintAllocation(ctx, args)
	return err
}
