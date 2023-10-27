package querytest

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFetchAllocatedComplaints(t *testing.T) {
	allocated_id, err := uuid.Parse("129593c1-71b9-4e7e-9da0-a25518c4cd66")

	require.NoError(t, err)

	allocated_complaints, err := testQueries.FetchAllocatedComplaintByEmp(context.Background(), allocated_id)
	require.NoError(t, err)
	fmt.Println("LIST : ", len(allocated_complaints))
	require.NotEmpty(t, allocated_complaints)
}
