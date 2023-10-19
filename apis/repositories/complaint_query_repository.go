package repositories

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/aniket-skroman/skroman_support_installation/utils"
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

func (repo *complaint_repository) UpdateComplaintInfo(args db.UpdateComplaintInfoParams) (db.ComplaintInfo, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.UpdateComplaintInfo(ctx, args)
}

func (repo *complaint_repository) FetchDeviceFileById(file_id uuid.UUID) (db.DeviceImages, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.FetchDeviceFileById(ctx, file_id)
}

// delete a device file , image/mp4
func (repo *complaint_repository) DeleteDeviceFiles(file_id uuid.UUID) (sql.Result, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.DeleteDeviceFiles(ctx, file_id)
}

// delete a complaint , remove all the data from DB
func (repo *complaint_repository) DeleteComplaint(complaint_id uuid.UUID) ([]db.DeviceImages, error) {
	// make a db transaction
	tx, err := repo.db.DB_instatnce().Begin()

	if err != nil {
		return nil, err
	}

	qtx := repo.db.WithTx(tx)

	// fetch a complaint info and store a complaint info id to remove device images/videos
	complaint_info, err := qtx.FetchComplaintDetailByComplaint(context.Background(), complaint_id)

	if err != nil {
		return nil, err
	}

	// collect all device file images/videos
	device_files, err := qtx.FetchDeviceImagesByComplaintId(context.Background(), complaint_info.ComplaintInfoID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	// now delete all device files
	for i := range device_files {
		_, err := qtx.DeleteDeviceFiles(context.Background(), device_files[i].ID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// delete complaint info by comaplint id
	delete_result, err := qtx.DeleteComplaintInfoBYId(context.Background(), complaint_id)
	if err != nil {
		return nil, err
	}

	affected_rows, _ := delete_result.RowsAffected()
	if affected_rows == 0 {
		return nil, errors.New(utils.DELETE_FAILED)
	}

	// delete complaint allocation
	_, err = qtx.DeleteComplaintAllcation(context.Background(), complaint_id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	// in last delete a orignal complaint
	_, err = qtx.DeleteComplaintByID(context.Background(), complaint_id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	return device_files, tx.Commit()
}
