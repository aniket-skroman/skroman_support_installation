package querytest

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFetchAllocatedComplaints(t *testing.T) {
	allocated_id, err := uuid.Parse("e73c4be4-ba03-49d0-8848-107a76802202")

	require.NoError(t, err)
	arg := db.FetchAllocatedComplaintsByEmpPendingParams{
		AllocatedTo: allocated_id,
		CreatedAt:   time.Now().AddDate(0, 0, -1),
	}
	allocated_complaints, err := testQueries.FetchAllocatedComplaintsByEmpPending(context.Background(), arg)
	require.NoError(t, err)
	fmt.Println("LIST : ", len(allocated_complaints))
	for i := range allocated_complaints {
		fmt.Println("List Data : ", allocated_complaints[i].OnDate)
	}
	require.NotEmpty(t, allocated_complaints)
}

func createRandomProgress(t *testing.T) {

	complaint_id, _ := uuid.Parse("0a5b9181-9326-40c4-8efe-15542d86f458")

	args := db.CreateComplaintProgressParams{
		ComplaintID:       complaint_id,
		ProgressStatement: "test statement",
		StatementBy:       uuid.New(),
	}

	complaint_progess, err := testQueries.CreateComplaintProgress(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, complaint_progess)

}

func TestCreateComplaintProgress(t *testing.T) {
	createRandomProgress(t)
}

func TestFetchComplaintProgess(t *testing.T) {
	complaint_id, _ := uuid.Parse("2661eccf-6dbd-4927-9c1b-e6b3316713f7")

	complaints, err := testQueries.FetchComplaintProgress(context.Background(), complaint_id)

	require.NoError(t, err)
	require.NotEmpty(t, complaints)
}

func TestDeleteComplaintProgress(t *testing.T) {

	id, _ := uuid.Parse("40501426-a017-4845-bd90-d9a553fdaeb8")

	result, err := testQueries.DeleteComplaintProgressById(context.Background(), id)

	require.NoError(t, err)
	affected_rows, _ := result.RowsAffected()
	require.NotZero(t, affected_rows)
}
