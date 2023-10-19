package repositories

import (
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

func (repo *allocation_repository) CreateComplaintAllocation(args db.CreateComplaintAllocationParams) error {
	ctx, cancel := repo.Init()
	defer cancel()

	_, err := repo.db.Queries.CreateComplaintAllocation(ctx, args)
	return err
}

func (repo *allocation_repository) UpdateComplaintAllocation(args db.UpdateComplaintAllocationParams) error {
	ctx, cancel := repo.Init()
	defer cancel()
	_, err := repo.db.Queries.UpdateComplaintAllocation(ctx, args)
	return err
}

func (repo *allocation_repository) FetchAllocationByComplaintId(complaint_id uuid.UUID) (db.ComplaintAllocations, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.FetchComplaintAllocationByComplaint(ctx, complaint_id)
}
