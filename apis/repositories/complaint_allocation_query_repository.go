package repositories

import (
	"database/sql"
	"errors"

	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

func (repo *allocation_repository) CreateComplaintAllocation(args db.CreateComplaintAllocationParams) error {
	// implement a DB transaction
	tx, err := repo.db.DB_instatnce().Begin()

	if err != nil {
		return err
	}

	qtx := repo.db.WithTx(tx)

	ctx, cancel := repo.Init()
	defer cancel()

	// create a complaint allocation first
	_, err = qtx.CreateComplaintAllocation(ctx, args)

	if err != nil {
		return err
	}

	// update complaint status in complaint info
	complaint_indo := db.UpdateComplaintStatusParams{
		Status:      "ALLOCATE",
		ComplaintID: args.ComplaintID,
	}

	result, err := qtx.UpdateComplaintStatus(ctx, complaint_indo)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	afftected_rows, _ := result.RowsAffected()
	if afftected_rows == 0 {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return errors.New("failed to allocate a complaint")
	}
	err = tx.Commit()
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

// check for duplicate complaint allocation
func (repo *allocation_repository) CheckDuplicateComplaintAllocation(complaint_id uuid.UUID) (sql.Result, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CheckDuplicateComplaintAllocation(ctx, complaint_id)
}

// check complaint status before allocate emp or update emp, complaint status should be init/allocate should not be complete
func (repo *allocation_repository) CheckComplaintStatusBeforeUpdate(complaint_id uuid.UUID) (sql.Result, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CheckComplaintStatusBeforeUpdate(ctx, complaint_id)
}
