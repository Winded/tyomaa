package api

type ClockGetResponse struct {
	Entry *TimeEntry `json:"entry"`
}

type ClockStartPostRequest struct {
	Project string `json:"project"`
}
type ClockStartPostResponse struct {
	Entry TimeEntry `json:"entry"`
}

type ClockStopPostResponse ClockStartPostResponse
