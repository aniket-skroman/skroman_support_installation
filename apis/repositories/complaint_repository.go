package repositories

import (
	"context"
	"time"

	"github.com/aniket-skroman/skroman_support_installation/apis"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
)

type ComplaintRepository interface {
	Init() (context.Context, context.CancelFunc)
	CreateComplaint(db.CreateComplaintParams) (db.Complaints, error)
	CreateComplaintInfo(db.CreateComplaintInfoParams) (db.ComplaintInfo, error)
	FetchAllComplaints(db.FetchAllComplaintsParams) ([]db.ComplaintInfo, error)
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
