package api

type EntriesGetRequest struct {
	// TODO
}
type EntriesGetResponse struct {
	Entries []TimeEntry `json:"entries"`
}

type EntriesPostRequest struct {
	Entry TimeEntry `json:"entry"`
}

type EntriesPostResponse struct {
	Entry TimeEntry `json:"entry"`
}

type EntriesSingleGetResponse struct {
	Entry TimeEntry `json:"entry"`
}

type EntriesSinglePostRequest EntriesPostRequest
type EntriesSinglePostResponse EntriesPostResponse
