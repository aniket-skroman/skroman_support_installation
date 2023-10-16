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
	FetchAllComplaints(db.FetchAllComplaintsParams) ([]db.ComplaintInfo, error)
	CountComplaints() (sql.Result, error)
	FetchComplaintDetailByComplaint(uuid.UUID) (db.FetchComplaintDetailByComplaintRow, error)
	FetchDeviceImagesByComplaintId(uuid.UUID) ([]db.DeviceImages, error)
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
