package repositories

import (
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
)

func (repo *complaint_repository) CreateComplaint(args db.CreateComplaintParams) (db.Complaints, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CreateComplaint(ctx, args)
}

func (repo *complaint_repository) CreateComplaintInfo(args db.CreateComplaintInfoParams) (db.ComplaintInfo, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CreateComplaintInfo(ctx, args)
}

func (repo *complaint_repository) FetchAllComplaints(args db.FetchAllComplaintsParams) ([]db.ComplaintInfo, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.FetchAllComplaints(ctx, args)
}
