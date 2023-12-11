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
