package proxycalls

import (
	"bytes"
	"net/http"
)

var (
	BASE_URL     = "http://3.7.18.55:3000/skroman/"
	USER_SERVICE = "http://15.207.19.172:8080/api/"
	// USER_SERVICE = "http://localhost:8080/api/"
)

type APIRequest interface {
	MakeApiRequest() (*http.Response, error)
}

type ProxyCalls struct {
	ReqEndpoint   string
	RequestBody   []byte
	RequestMethod string
	RequestParams interface{}
	RequestHeader map[string]string
	IsRequestBody bool // request need a body or params
	Response      *http.Response
}

func NewAPIRequest(endpoint, method string, isRequestBody bool, requestBody []byte, params interface{}, headers map[string]string) APIRequest {
	return &ProxyCalls{
		ReqEndpoint:   USER_SERVICE + endpoint,
		RequestMethod: method,
		RequestParams: params,
		RequestHeader: headers,
		IsRequestBody: isRequestBody,
	}
}

func (pc *ProxyCalls) MakeApiRequest() (*http.Response, error) {
	var request *http.Request
	var err error

	if pc.IsRequestBody {
		request, err = http.NewRequest(pc.RequestMethod, pc.ReqEndpoint, bytes.NewReader(pc.RequestBody))
	} else {
		request, err = http.NewRequest(pc.RequestMethod, pc.ReqEndpoint, nil)
	}

	if err != nil {
		return nil, err
	}

	for key, val := range pc.RequestHeader {
		request.Header.Set(key, val)
	}
	request.Close = true
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)

	return response, err
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
