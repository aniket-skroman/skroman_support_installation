package services

import (
	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/helper"
	"github.com/aniket-skroman/skroman_support_installation/apis/repositories"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
)

type ComplaintAllocationService interface {
	AllocateComplaint(req dto.CreateAllocationRequestDTO) error
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
