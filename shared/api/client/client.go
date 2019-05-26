package client

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/winded/tyomaa/shared/api"
)

func makeError(response *Response) (err api.ApiError) {
	err.Status = response.StatusCode
	json.Unmarshal(response.Body, &err)
	return err
}

type Settings struct {
	Host  string `json:"host"`
	Token string `json:"token"`
}

type ApiClient struct {
	Settings Settings
	client   IHttpClient
}

func NewApiClient(httpClient IHttpClient, settings Settings) *ApiClient {
	return &ApiClient{
		Settings: settings,
		client:   httpClient,
	}
}

func (this *ApiClient) headers() map[string]string {
	return map[string]string{
		"Content-Type":   "application/json",
		"X-Access-Token": this.Settings.Token,
	}
}

func (this *ApiClient) do(method, path string, query url.Values, inputData, outputData interface{}) (err error) {
	sq := query.Encode()
	if sq != "" {
		sq = "?" + sq
	}

	endpoint := this.Settings.Host + path + sq

	var reqBody []byte
	if inputData != nil {
		reqBody, err = json.Marshal(inputData)
		if err != nil {
			return
		}
	}

	resp, err := this.client.Do(&Request{
		Method: method,
		URL:    endpoint,
		Header: this.headers(),
		Body:   reqBody,
	})
	if err != nil {
		return
	}

	if resp.StatusCode >= api.StatusBadRequest {
		return makeError(resp)
	}

	err = json.Unmarshal(resp.Body, outputData)
	return
}

func (this *ApiClient) AuthTokenGet() (response api.TokenGetResponse, err error) {
	err = this.do(api.MethodGet, "/auth/token", url.Values{}, nil, &response)
	return
}
func (this *ApiClient) AuthTokenPost(request api.TokenPostRequest) (response api.TokenPostResponse, err error) {
	err = this.do(api.MethodPost, "/auth/token", url.Values{}, &request, &response)
	return
}

func (this *ApiClient) ClockGet() (response api.ClockGetResponse, err error) {
	err = this.do(api.MethodGet, "/clock", url.Values{}, nil, &response)
	return
}

func (this *ApiClient) ClockStartPost(request api.ClockStartPostRequest) (response api.ClockStartPostResponse, err error) {
	err = this.do(api.MethodPost, "/clock/start", url.Values{}, &request, &response)
	return
}

func (this *ApiClient) ClockStopPost() (response api.ClockStopPostResponse, err error) {
	err = this.do(api.MethodPost, "/clock/stop", url.Values{}, nil, &response)
	return
}

func (this *ApiClient) EntriesGet(request api.EntriesGetRequest) (response api.EntriesGetResponse, err error) {
	err = this.do(api.MethodGet, "/entries", url.Values{}, &request, &response)
	return
}
func (this *ApiClient) EntriesPost(request api.EntriesPostRequest) (response api.EntriesPostResponse, err error) {
	err = this.do(api.MethodPost, "/entries", url.Values{}, &request, &response)
	return
}

func (this *ApiClient) EntriesSingleGet(entryId uint) (response api.EntriesSingleGetResponse, err error) {
	err = this.do(api.MethodGet, "/entries/"+strconv.Itoa(int(entryId)), url.Values{}, nil, &response)
	return
}
func (this *ApiClient) EntriesSinglePost(entryId uint, request api.EntriesSinglePostRequest) (response api.EntriesSinglePostResponse, err error) {
	err = this.do(api.MethodPost, "/entries/"+strconv.Itoa(int(entryId)), url.Values{}, &request, &response)
	return
}
func (c *ApiClient) EntriesSingleDelete(entryId uint) (response api.EntriesSingleDeleteResponse, err error) {
	err = c.do(api.MethodDelete, "/entries/"+strconv.Itoa(int(entryId)), url.Values{}, nil, &response)
	return
}

func (this *ApiClient) ProjectsGet() (response api.ProjectsGetResponse, err error) {
	err = this.do(api.MethodGet, "/projects", url.Values{}, nil, &response)
	return
}
