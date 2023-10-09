package querytest

import (
	"context"
	"database/sql"
	"testing"
	"time"

	sqlc_lib "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func create_complaint(t *testing.T) sqlc_lib.Complaint {
	args := sqlc_lib.CreateComplaintParams{
		ClientID:         uuid.New(),
		DeviceID:         uuid.New(),
		ProblemStatement: utils.RandomString(15),
		ProblemCategory: sql.NullString{
			String: "TEST",
			Valid:  true,
		},
		ClientAvailable: time.Now(),
		Status:          "INIT",
	}

	complaint, err := testQueries.CreateComplaint(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, complaint)

	require.Equal(t, args.ClientID, complaint.ClientID)

	return complaint
}

func TestCreateComplaint(t *testing.T) {
	create_complaint(t)
}

func TestFetchAllComplaints(t *testing.T) {
	args := sqlc_lib.GetComplaintsParams{
		Limit:  10,
		Offset: 0,
	}

	complaints, err := testQueries.GetComplaints(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, complaints)

}
