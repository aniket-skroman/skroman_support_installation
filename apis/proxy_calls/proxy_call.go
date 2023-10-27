package proxycalls

import (
	"bytes"
	"net/http"
)

var (
	BASE_URL     = "http://3.7.18.55:3000/skroman/"
	USER_SERVICE = "http://15.207.19.172:8080/api/"
)

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
