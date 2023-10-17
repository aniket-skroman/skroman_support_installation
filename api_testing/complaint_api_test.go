package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"

	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/stretchr/testify/require"
)

var debug_logger *log.Logger
var Token string

func init() {

	debug_logger = log.New(os.Stdout, "DEBUG : ", log.Flags())
	Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmEwZGFhYjItZDY1MS00MmIxLThjOWYtOWY5ODcxMTE4YzdlIiwidXNlcl90eXBlIjoiQURNSU4iLCJjcmVhdGVkX2F0IjoiMjAyMy0xMC0xN1QwNToxMzo0OC45NDI4NDQwNzFaIiwiZXhwIjozMDAwMDAwMDAwMDAsImlhdCI6MTY5NzUxOTAyOCwiaXNzIjoieWRobndiIn0.TCXCBCz9xRig5KKT4waRDr1NV-6Qg5ETE-6Irln-0M0"
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

func TestValidateToken(t *testing.T) {

	requrl := "http://localhost:8080/api/fetch-user"

	request, err := http.NewRequest(http.MethodGet, requrl, nil)

	require.NoError(t, err)
	request.Header.Set("Authorization", Token)

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	require.NotEmpty(t, response)

	response_body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	fmt.Println("RESPONSE BODY : ", string(response_body))
	require.Equal(t, response.StatusCode, http.StatusOK)
}
