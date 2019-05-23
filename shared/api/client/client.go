package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/winded/tyomaa/shared/api"
)

func makeError(response *http.Response) (err api.ApiError) {
	err.Status = response.StatusCode

	data, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return err
	}

	json.Unmarshal(data, &err)
	return err
}

type Settings struct {
	Host  string `json:"host"`
	Token string `json:"token"`
}

type ApiClient struct {
	Settings Settings
	client   http.Client
}

func NewApiClient(settings Settings) *ApiClient {
	return &ApiClient{
		Settings: settings,
		client: http.Client{
			Timeout: time.Second * 20,
		},
	}
}

func (this *ApiClient) setHeaders(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Access-Token", this.Settings.Token)
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

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}
	this.setHeaders(req)

	resp, err := this.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return makeError(resp)
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(resBody, outputData)
	return
}

func (this *ApiClient) AuthTokenGet() (response api.TokenGetResponse, err error) {
	err = this.do(http.MethodGet, "/auth/token", url.Values{}, nil, &response)
	return
}
func (this *ApiClient) AuthTokenPost(request api.TokenPostRequest) (response api.TokenPostResponse, err error) {
	err = this.do(http.MethodPost, "/auth/token", url.Values{}, &request, &response)
	return
}

func (this *ApiClient) ClockGet() (response api.ClockGetResponse, err error) {
	err = this.do(http.MethodGet, "/clock", url.Values{}, nil, &response)
	return
}

func (this *ApiClient) ClockStartPost(request api.ClockStartPostRequest) (response api.ClockStartPostResponse, err error) {
	err = this.do(http.MethodPost, "/clock/start", url.Values{}, &request, &response)
	return
}

func (this *ApiClient) ClockStopPost() (response api.ClockStopPostResponse, err error) {
	err = this.do(http.MethodPost, "/clock/stop", url.Values{}, nil, &response)
	return
}

func (this *ApiClient) EntriesGet(request api.EntriesGetRequest) (response api.EntriesGetResponse, err error) {
	err = this.do(http.MethodGet, "/entries", url.Values{}, &request, &response)
	return
}
func (this *ApiClient) EntriesPost(request api.EntriesPostRequest) (response api.EntriesPostResponse, err error) {
	err = this.do(http.MethodPost, "/entries", url.Values{}, &request, &response)
	return
}

func (this *ApiClient) EntriesSingleGet(entryId uint) (response api.EntriesSingleGetResponse, err error) {
	err = this.do(http.MethodGet, "/entries/"+strconv.Itoa(int(entryId)), url.Values{}, nil, &response)
	return
}
func (this *ApiClient) EntriesSinglePost(entryId uint, request api.EntriesSinglePostRequest) (response api.EntriesSinglePostResponse, err error) {
	err = this.do(http.MethodPost, "/entries/"+strconv.Itoa(int(entryId)), url.Values{}, &request, &response)
	return
}

func (this *ApiClient) ProjectsGet() (response api.ProjectsGetResponse, err error) {
	err = this.do(http.MethodGet, "/projects", url.Values{}, nil, &response)
	return
}
