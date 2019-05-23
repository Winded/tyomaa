package api

import "time"

type Project struct {
	Name      string        `json:"name"`
	TotalTime time.Duration `json:"totalTime"`
}

type ProjectsGetResponse struct {
	Projects []Project `json:"projects"`
}
