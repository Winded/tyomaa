package xhr_http_client

import (
	"time"

	"honnef.co/go/js/xhr"

	"github.com/winded/tyomaa/shared/api/client"
)

// XHRHttpCLient is an IHttpClient implementation for the API that directly uses XHR through honnef.co/go/js/xhr.
// This is because standard net/http causes issue due to a compilation error in GopherJS. See: https://github.com/gopherjs/gopherjs/issues/604
type XHRHttpCLient struct {
	timeout time.Duration
}

func New(timeout time.Duration) *XHRHttpCLient {
	return &XHRHttpCLient{
		timeout: timeout,
	}
}

func (this *XHRHttpCLient) Do(request *client.Request) (*client.Response, error) {
	req := xhr.NewRequest(request.Method, request.URL)
	req.Timeout = int(this.timeout / time.Millisecond)
	req.ResponseType = "application/json"

	for header, value := range request.Header {
		req.SetRequestHeader(header, value)
	}

	err := req.Send(request.Body)
	if err != nil {
		return nil, err
	}

	return &client.Response{
		StatusCode: req.Status,
		Body:       []byte(req.ResponseText),
	}, nil
}
