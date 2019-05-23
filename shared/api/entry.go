package api

import "time"

type TimeEntry struct {
	ID      uint       `json:"id"`
	Project string     `json:"project"`
	Start   time.Time  `json:"start"`
	End     *time.Time `json:"end"`
}
