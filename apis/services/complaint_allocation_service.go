package services

import (
	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/helper"
	"github.com/aniket-skroman/skroman_support_installation/apis/repositories"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

type ComplaintAllocationService interface {
	AllocateComplaint(req dto.CreateAllocationRequestDTO) error
	UpdateComplaintAllocation(req dto.UpdateAllocateComplaintRequestDTO) error
	FetchAllocationByComplaintId(complaint_id string) (db.ComplaintAllocations, error)
}

type allocation_service struct {
	allocation_repo repositories.ComplaintAllocationRepository
}

func NewComplaintAllocationService(repo repositories.ComplaintAllocationRepository) ComplaintAllocationService {
	return &allocation_service{
		allocation_repo: repo,
	}
}

func (ser *allocation_service) AllocateComplaint(req dto.CreateAllocationRequestDTO) error {
	complaint_obj, err := helper.ValidateUUID(req.ComplaintId)
	if err != nil {
		return err
	}

	allocate_by, err := helper.ValidateUUID(req.AllocateBy)
	if err != nil {
		return err
	}

	allocate_to, err := helper.ValidateUUID(req.AllocateTo)
	if err != nil {
		return err
	}

	args := db.CreateComplaintAllocationParams{
		ComplaintID: complaint_obj,
		AllocatedBy: allocate_by,
		AllocatedTo: allocate_to,
	}

	err = ser.allocation_repo.CreateComplaintAllocation(args)
	err = helper.Handle_db_err(err)
	return err
}

func (ser *allocation_service) UpdateComplaintAllocation(req dto.UpdateAllocateComplaintRequestDTO) error {
	id, err := helper.ValidateUUID(req.Id)

	if err != nil {
		return err
	}
	allocate_to, err := helper.ValidateUUID(req.AllocateTo)
	if err != nil {
		return err
	}

	allocate_by, err := helper.ValidateUUID(req.AllocateBy)
	if err != nil {
		return err
	}

	args := db.UpdateComplaintAllocationParams{
		ID:          id,
		AllocatedTo: allocate_to,
		AllocatedBy: allocate_by,
	}

	err = ser.allocation_repo.UpdateComplaintAllocation(args)

	if err != nil {
		return err
	}

	err = helper.Handle_db_err(err)
	return err
}

func (ser *allocation_service) FetchAllocationByComplaintId(complaint_id string) (db.ComplaintAllocations, error) {
	complaint_obj_id, err := uuid.Parse(complaint_id)

	if err != nil {
		return db.ComplaintAllocations{}, nil
	}

	return ser.allocation_repo.FetchAllocationByComplaintId(complaint_obj_id)
}
