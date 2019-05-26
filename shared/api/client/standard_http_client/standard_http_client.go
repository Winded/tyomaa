package standard_http_client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/winded/tyomaa/shared/api/client"
)

type StandardHttpClient struct {
	client http.Client
}

func New(timeout time.Duration) *StandardHttpClient {
	return &StandardHttpClient{
		client: http.Client{
			Timeout: timeout,
		},
	}
}

func (this *StandardHttpClient) Do(request *client.Request) (*client.Response, error) {
	req, err := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, err
	}

	for header, value := range request.Header {
		req.Header.Set(header, value)
	}

	resp, err := this.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &client.Response{
		StatusCode: resp.StatusCode,
		Body:       resBody,
	}, nil
}
