package apitesting

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/stretchr/testify/require"
)

var debug_logger *log.Logger
var Token string

func init() {
	log_file, err := os.Create("app.log")

	if err != nil {
		log.Fatal(err)
	}

	debug_logger = log.New(log_file, "DEBUG : ", log.Flags())
	Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDRiYTc3ZjMtZmJjNC00MjFiLTliZGEtMGQ5NTk4YWNlODMxIiwidXNlcl90eXBlIjoiZW1wIiwiY3JlYXRlZF9hdCI6IjIwMjMtMTAtMTNUMTA6Mzk6MDUuODA3MjcyMDY3WiIsImV4cCI6MzAwMDAwMDAwMDAwLCJpYXQiOjE2OTcxOTI5NDUsImlzcyI6InlkaG53YiJ9.23TzhMNfKUyWgvyC--wbcZA_iOVLax7fn5Qdk1oWCXU"
}

func TestCreateComplaint(t *testing.T) {

	args := []struct {
		TestName       string
		RequestBody    dto.CreateComplaintRequestDTO
		ExpectedErr    bool
		ExpectedStatus int
	}{
		{
			TestName: "FIRST",
			RequestBody: dto.CreateComplaintRequestDTO{
				ClientID:         "TEST-CLIENT-FIRST",
				DeviceID:         "TEST-DEVICE-FIRST",
				DeviceType:       "TEST-TYPE",
				DeviceModel:      "TEST-MODEL",
				ProblemStatement: utils.RandomString(13),
				ProblemCategory:  "TEST-CATEGORY",
				ClientAvailable:  "2023-10-13 15:36:03",
			},
			ExpectedErr:    false,
			ExpectedStatus: http.StatusOK,
		}, {
			TestName: "SECOND-Device-Model-Empty",
			RequestBody: dto.CreateComplaintRequestDTO{
				ClientID:         "TEST-CLIENT-FIRST",
				DeviceID:         "TEST-DEVICE-FIRST",
				DeviceType:       "TEST-TYPE",
				DeviceModel:      "",
				ProblemStatement: utils.RandomString(13),
				ProblemCategory:  "TEST-CATEGORY",
				ClientAvailable:  "2023-10-13 15:36:03",
			},
			ExpectedErr:    true,
			ExpectedStatus: http.StatusBadRequest,
		}, {
			TestName: "THIRD-Device-Type-Empty",
			RequestBody: dto.CreateComplaintRequestDTO{
				ClientID:         "TEST-CLIENT-FIRST",
				DeviceID:         "TEST-DEVICE-FIRST",
				DeviceType:       "",
				DeviceModel:      "TEST-MODEL",
				ProblemStatement: utils.RandomString(13),
				ProblemCategory:  "TEST-CATEGORY",
				ClientAvailable:  "2023-10-13 15:36:03",
			},
			ExpectedErr:    true,
			ExpectedStatus: http.StatusBadRequest,
		}, {
			TestName: "FOUTH-PRO-STATE-EMPTY",
			RequestBody: dto.CreateComplaintRequestDTO{
				ClientID:         "TEST-CLIENT-FIRST",
				DeviceID:         "TEST-DEVICE-FIRST",
				DeviceType:       "TEST-TYPE",
				DeviceModel:      "TEST-MODEL",
				ProblemStatement: "",
				ProblemCategory:  "TEST-CATEGORY",
				ClientAvailable:  "2023-10-13 15:36:03",
			},
			ExpectedErr:    true,
			ExpectedStatus: http.StatusBadRequest,
		}, {
			TestName: "FIFTH-PRO-CAT-EMPTY",
			RequestBody: dto.CreateComplaintRequestDTO{
				ClientID:         "TEST-CLIENT-FIRST",
				DeviceID:         "TEST-DEVICE-FIRST",
				DeviceType:       "TEST-TYPE",
				DeviceModel:      "TEST-MODEL",
				ProblemStatement: utils.RandomString(12),
				ProblemCategory:  "",
				ClientAvailable:  "2023-10-13 15:36:03",
			},
			ExpectedErr:    true,
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	url := "http://13.233.196.149:8181/api/create-complaint"

	for _, arg := range args {
		t.Run(arg.TestName, func(t *testing.T) {
			request_body, err := json.Marshal(arg.RequestBody)

			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(request_body))
			request.Header.Set("Authorization", Token)
			require.NoError(t, err)

			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err)
			response_body, err := io.ReadAll(response.Body)

			require.NoError(t, err)

			debug_logger.Println("REQUEST : ", request)
			debug_logger.Println("RESPONSE : ", string(response_body))
			debug_logger.Println("RESPONCE STATUS CODE : ", response.StatusCode)
			debug_logger.Println("EXPECTED CODE : ", arg.ExpectedStatus)
			debug_logger.Println()

			require.NoError(t, err)

			require.Equal(t, arg.ExpectedStatus, response.StatusCode)

		})
	}
}

func TestFetchComplaints(t *testing.T) {
	args := []struct {
		TestName       string
		Params         db.FetchAllComplaintsParams
		AccessToken    string
		ExpectedStatus int
	}{
		{
			TestName: "FIRST",
			Params: db.FetchAllComplaintsParams{
				Limit:  10,
				Offset: 1,
			},
			AccessToken:    Token,
			ExpectedStatus: http.StatusOK,
		}, {
			TestName: "SECOND",
			Params: db.FetchAllComplaintsParams{
				Limit:  10,
				Offset: 101,
			},
			AccessToken:    Token,
			ExpectedStatus: http.StatusNotFound,
		}, {
			TestName: "THIRD",
			Params: db.FetchAllComplaintsParams{
				Limit:  10,
				Offset: -1,
			},
			AccessToken:    Token,
			ExpectedStatus: http.StatusInternalServerError,
		},
		{
			TestName: "FOURTH",
			Params: db.FetchAllComplaintsParams{
				Limit:  10,
				Offset: 1,
			},
			AccessToken:    "",
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	url := "http://13.233.196.149:8181/api/fetch-complaints"

	for _, arg := range args {
		t.Run(arg.TestName, func(t *testing.T) {
			req_url := url + "/" + strconv.Itoa(int(arg.Params.Offset)) + "/" + strconv.Itoa(int(arg.Params.Limit))
			request, err := http.NewRequest(http.MethodGet, req_url, nil)

			require.NoError(t, err)
			request.Header.Set("Authorization", arg.AccessToken)
			// q := request.URL.Query()
			// q.Add("page_id", strconv.Itoa(int(arg.Params.Offset)))
			// q.Add("page_size", strconv.Itoa(int(arg.Params.Limit)))

			// request.URL.RawQuery = q.Encode()

			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err)

			response_body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			debug_logger.Println("REQUEST : ", request)
			debug_logger.Println("RESPONSE : ", string(response_body))
			debug_logger.Println("RESPONCE STATUS CODE : ", response.StatusCode)
			debug_logger.Println("EXPECTED CODE : ", arg.ExpectedStatus)
			debug_logger.Println()

			require.Equal(t, arg.ExpectedStatus, response.StatusCode)
		})
	}
}
