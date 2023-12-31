package proxycalls

import (
	"bytes"
	"log"
	"net/http"
	"os"
)

var (
	BASE_URL     = "http://3.7.18.55:3000/skroman/"
	USER_SERVICE = "http://localhost:8080/api/"
	logger       *log.Logger
)

func init() {
	file, err := os.Open("./api_testing/app.log")

	if err != nil {
		log.Fatal(err)
	}

	logger = log.New(file, "DEBUG : ", log.Flags())
}

type ProxyCalls struct {
	ReqEndpoint   string
	RequestBody   []byte
	RequestMethod string
	RequestParams interface{}
	IsRequestBody bool // request need a body or params
	Response      *http.Response
}

func (procall *ProxyCalls) MakeRequestWithBody() (*http.Response, error) {
	request, err := http.NewRequest(procall.RequestMethod, BASE_URL+procall.ReqEndpoint, bytes.NewReader(procall.RequestBody))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	return response, err
}
