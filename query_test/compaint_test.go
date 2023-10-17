package querytest

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

var aws_access_key_id = "AKIA3VMV3LWIQ6EL63WU"
var aws_secret_access_key = "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"
var region = "ap-south-1"

var bucket_name = "skromansupportbucket"

// return file path
func upload_image() string {
	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
	_, err := creds.Get()
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	cfg := aws.NewConfig().WithRegion("ap-south-1").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	file, err := os.Open("./test.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := "media/" + file.Name()
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket_name),
		Key:    aws.String(path),
		Body:   fileBytes,

		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	fmt.Printf("response %s", awsutil.StringValue(resp))
	return path
}

func TestUploadDeviceImages(t *testing.T) {
	path := upload_image()
	complaint_id, err := uuid.Parse("3f5263f3-897f-463e-aa8a-186ab98ef371")

	require.NoError(t, err)

	args := db.UploadDeviceImagesParams{
		ComplaintInfoID: complaint_id,
		DeviceImage:     path,
	}

	device, err := testQueries.UploadDeviceImages(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, device)

	fmt.Printf("%v+\n", device)
}

func TestFetchImageFromS3(t *testing.T) {
	complaint_id, err := uuid.Parse("3f5263f3-897f-463e-aa8a-186ab98ef371")
	require.NoError(t, err)

	device_images, err := testQueries.FetchDeviceImagesByComplaintId(context.Background(), complaint_id)

	require.NoError(t, err)
	require.NotEmpty(t, device_images)
	url := "https://skromansupportbucket.s3.ap-south-1.amazonaws.com/"
	for _, img := range device_images {
		img_path := fmt.Sprintf("%s%s", url, img.DeviceImage)

		fmt.Println(img_path)
	}
}

func TestCountComplaint(t *testing.T) {
	result, err := testQueries.CountComplaints(context.Background())

	require.NoError(t, err)

	affetcted_rows, err := result.RowsAffected()
	fmt.Println("Affected rows : ", affetcted_rows)
	require.NoError(t, err)
	require.NotZero(t, affetcted_rows)
}

func TestFetchComplaintByComplaintID(t *testing.T) {
	complaint_id, err := uuid.Parse("2571e975-4c0d-42fd-bea2-bdf59b38d482")
	require.NoError(t, err)

	complaints, err := testQueries.FetchComplaintDetailByComplaint(context.Background(), complaint_id)
	fmt.Printf("%+v\n", complaints)

	require.NoError(t, err)
	require.NotEmpty(t, complaints)

}

func TestProxyAPI(t *testing.T) {
	reqUrl := "http://3.7.18.55:3000/skroman/profileapi/profileuser/userId"
	body := struct {
		UserId string `json:"userId"`
	}{UserId: "User_id-iYfdKPhPS"}

	request_body, err := json.Marshal(&body)

	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewReader(request_body))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	// req_body, _ := io.ReadAll(request.Body)
	// fmt.Println("Request data : \n", string(req_body))
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	response_body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	contain_data := string(response_body)
	fmt.Println("Contain Data : \n", contain_data)
	require.NotEmpty(t, contain_data)
}
