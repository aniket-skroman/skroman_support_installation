package querytest

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CreateComplaint struct {
	Complaint     db.Complaints
	ComplaintInfo db.ComplaintInfo
	DeviceImages  db.DeviceImages
	ExpectErr     bool
}

func generate_random_complaint() db.Complaints {
	return db.Complaints{
		CreatedBy: uuid.New(),
	}
}

func generate_random_complaint_info(complaint_id uuid.UUID) db.ComplaintInfo {
	return db.ComplaintInfo{
		ComplaintID:      complaint_id,
		DeviceID:         "dk-device",
		ProblemStatement: utils.RandomString(15),
		ProblemCategory: sql.NullString{
			String: "TEST-CAT",
			Valid:  true,
		},
		ClientAvailable: time.Now().Add(15 * time.Hour),
		Status:          "INIT",
	}
}

func generate_random_device_imges(complaint_info_id uuid.UUID) db.DeviceImages {
	return db.DeviceImages{
		ComplaintInfoID: complaint_info_id,
		DeviceImage:     fmt.Sprintf("%s.jpg", complaint_info_id.String()),
	}
}

func make_transaction2(t *testing.T) {
	tests := []CreateComplaint{}

	for i := 0; i < 5; i++ {
		complaint := generate_random_complaint()
		complaint_info := generate_random_complaint_info(complaint.ID)
		device_img := generate_random_device_imges(complaint_info.ID)
		new_complaint := CreateComplaint{
			Complaint:     complaint,
			ComplaintInfo: complaint_info,
			DeviceImages:  device_img,
		}

		if i == 4 {
			device_img.DeviceImage = ""
			new_complaint.ExpectErr = true
		} else {
			new_complaint.ExpectErr = false
		}

		tests = append(tests, new_complaint)
	}

	for _, tt := range tests {
		t.Run(string(tt.Complaint.CreatedBy.String()), func(t *testing.T) {
			tx, err := testDB.Begin()
			assert.NoError(t, err)

			qxt := testQueries.WithTx(tx)

			//create complaint
			comp_args := db.CreateComplaintParams{
				ClientID:  "Test Clinet",
				CreatedBy: uuid.New(),
			}
			complaint, err := qxt.CreateComplaint(context.Background(), comp_args)

			assert.NoError(t, err)
			assert.NotEmpty(t, complaint)

			// create complaint info against complaint
			args := db.CreateComplaintInfoParams{
				ComplaintID:      complaint.ID,
				DeviceID:         tt.ComplaintInfo.DeviceID,
				ProblemStatement: tt.ComplaintInfo.ProblemStatement,
				ProblemCategory:  tt.ComplaintInfo.ProblemCategory,
				ClientAvailable:  tt.ComplaintInfo.ClientAvailable,
				Status:           tt.ComplaintInfo.Status,
			}

			complaint_info, err := qxt.CreateComplaintInfo(context.Background(), args)
			assert.NoError(t, err)
			assert.NotEmpty(t, complaint_info)
			assert.Equal(t, args.ComplaintID, complaint_info.ComplaintID)

			// add a device images, with actual issue
			device_args := db.AddDeviceImagesParams{
				ComplaintInfoID: complaint_info.ID,
				DeviceImage:     tt.DeviceImages.DeviceImage,
			}

			device_img, err := qxt.AddDeviceImages(context.Background(), device_args)
			assert.NoError(t, err)
			assert.NotEmpty(t, device_img)

			tx.Commit()
		})
	}
}

func TestCreateComplaint(t *testing.T) {
	make_transaction2(t)
}

func TestCreateDummyComplaint(t *testing.T) {
	args := db.CreateComplaintParams{
		ClientID:  "test-client",
		CreatedBy: uuid.New(),
	}
	complaint, err := testQueries.CreateComplaint(context.Background(), args)
	require.NoError(t, err)

	require.NotEmpty(t, complaint)
}

func TestCreateComplaintInfo(t *testing.T) {
	args := db.CreateComplaintParams{
		ClientID:  "test-client",
		CreatedBy: uuid.New(),
	}
	complaint, err := testQueries.CreateComplaint(context.Background(), args)
	require.NoError(t, err)

	require.NotEmpty(t, complaint)

	comp_info_args := db.CreateComplaintInfoParams{
		ComplaintID: complaint.ID,
		DeviceID:    "TEST-DEVICE",
		DeviceType: sql.NullString{
			String: "Test",
			Valid:  true,
		},
		DeviceModel: sql.NullString{
			String: "TEST-MODEL",
			Valid:  true,
		},
		ProblemStatement: "TEST-PROBLEM",
		ProblemCategory: sql.NullString{
			String: "TEST-Category",
			Valid:  true,
		},
		ClientAvailable: time.Now(),
	}

	comp_info, err := testQueries.CreateComplaintInfo(context.Background(), comp_info_args)
	require.NoError(t, err)
	require.NotEmpty(t, comp_info)
}

func TestFetchAllComplaints(t *testing.T) {

	args := []struct {
		TestName       string
		Params         db.FetchAllComplaintsParams
		ExpectedErr    bool
		ExpectedResult bool
	}{
		{
			TestName: "FIRST",
			Params: db.FetchAllComplaintsParams{
				Limit:  10,
				Offset: 1,
			},
			ExpectedErr:    false,
			ExpectedResult: true,
		}, {
			TestName: "SECOND",
			Params: db.FetchAllComplaintsParams{
				Limit:  30,
				Offset: -1,
			},
			ExpectedErr:    true,
			ExpectedResult: false,
		}, {
			TestName: "THIRD",
			Params: db.FetchAllComplaintsParams{
				Limit:  20,
				Offset: 16,
			},
			ExpectedErr:    false,
			ExpectedResult: false,
		},
	}

	for _, arg := range args {
		t.Run(arg.TestName, func(t *testing.T) {

			result, err := testQueries.FetchAllComplaints(context.Background(), arg.Params)

			if arg.ExpectedErr == true {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if arg.ExpectedResult == true {
				require.NotEmpty(t, result)
			} else {
				require.Empty(t, result)
			}
		})
	}
}
