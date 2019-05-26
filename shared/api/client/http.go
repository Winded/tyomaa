package client

type Request struct {
	Method string
	Header map[string]string
	URL    string
	Body   []byte
}

type Response struct {
	StatusCode int
	Body       []byte
}

// IHttpClient is a HTTP client abstraction for the API client.
// In some environments, like GopherJS, it's better to use other packages than standard net/http to make requests
type IHttpClient interface {
	Do(request *Request) (*Response, error)
}
