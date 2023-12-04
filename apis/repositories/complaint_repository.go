package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aniket-skroman/skroman_support_installation/apis"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
)

type ComplaintRepository interface {
	Init() (context.Context, context.CancelFunc)
	CreateComplaint(db.CreateComplaintParams) (db.Complaints, error)
	CreateComplaintInfo(db.CreateComplaintInfoParams) (db.ComplaintInfo, error)
	FetchAllComplaints(db.FetchAllComplaintsParams) ([]db.FetchAllComplaintsRow, error)
	FetchTotalComplaints(args db.FetchTotalComplaintsParams) ([]db.FetchTotalComplaintsRow, error)
	CountComplaints(status string) (sql.Result, error)
	FetchComplaintDetailByComplaint(uuid.UUID) (db.FetchComplaintDetailByComplaintRow, error)
	FetchDeviceImagesByComplaintId(uuid.UUID) ([]db.DeviceImages, error)
	UploadDeviceImage(db.UploadDeviceImagesParams) (db.DeviceImages, error)
	UpdateComplaintInfo(db.UpdateComplaintInfoParams) (db.ComplaintInfo, error)
	DeleteDeviceFiles(file_id uuid.UUID) (sql.Result, error)
	FetchDeviceFileById(file_id uuid.UUID) (db.DeviceImages, error)
	DeleteComplaint(complaint_id uuid.UUID) ([]db.DeviceImages, error)
	FetchComplaintById(complaint_id uuid.UUID) (db.FetchComplaintByComplaintIdRow, error)
	UpdateComplaintStatus(args db.UpdateComplaintStatusParams) (sql.Result, error)
	AllComplaintsCount() (db.CountAllComplaintRow, error)
	FetchCountByMonths() ([]db.FetchCountByMonthRow, error)
	FetchComplaintStatus(complaint_info_id uuid.UUID) (string, error)
	FetchComplaintsByClient(args db.FetchComplaintsByClientParams) ([]db.FetchComplaintsByClientRow, error)
	CountComplaintByClient(client_id string) (int64, error)
	TotalComplaints() (int64, error)
}

type complaint_repository struct {
	db *apis.Store
}

func NewComplaintRepository(db *apis.Store) ComplaintRepository {
	return &complaint_repository{
		db: db,
	}
}

func (repo *complaint_repository) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel
}
