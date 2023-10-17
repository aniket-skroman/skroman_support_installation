package repositories

import (
	"database/sql"

	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
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

func (repo *complaint_repository) CountComplaints() (sql.Result, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CountComplaints(ctx)
}

func (repo *complaint_repository) FetchComplaintDetailByComplaint(complaint_id uuid.UUID) (db.FetchComplaintDetailByComplaintRow, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.FetchComplaintDetailByComplaint(ctx, complaint_id)
}

func (repo *complaint_repository) FetchDeviceImagesByComplaintId(complaint_info_id uuid.UUID) ([]db.DeviceImages, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.FetchDeviceImagesByComplaintId(ctx, complaint_info_id)
}

func (repo *complaint_repository) UploadDeviceImage(args db.UploadDeviceImagesParams) (db.DeviceImages, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.UploadDeviceImages(ctx, args)
}
