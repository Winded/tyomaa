package models

import (
	"github.com/go-humble/locstor"
	"github.com/winded/tyomaa/frontend/js/dispatcher"
	"github.com/winded/tyomaa/shared/api"
	"github.com/winded/tyomaa/shared/api/client"
)

type Clock interface {
	ActiveProject() string
	ActiveEntry() (api.TimeEntry, error)

	Start(project string) (api.TimeEntry, error)
	Stop() (api.TimeEntry, error)
}

type clock struct {
	apiClient  *client.ApiClient
	dispatcher *dispatcher.Dispatcher

	activeEntry api.TimeEntry
}

func NewClock(apiClient *client.ApiClient, disp *dispatcher.Dispatcher) *clock {
	return &clock{
		apiClient:  apiClient,
		dispatcher: disp,
	}
}

func (c *clock) ActiveProject() string {
	project, _ := locstor.GetItem("active_project")
	return project
}

func (c *clock) ActiveEntry() (api.TimeEntry, error) {
	if c.activeEntry.ID != 0 {
		return c.activeEntry, nil
	}

	response, err := c.apiClient.ClockGet()
	if err != nil {
		return api.TimeEntry{}, err
	}

	c.activeEntry = response.Entry
	return c.activeEntry, nil
}

func (c *clock) Start(project string) (api.TimeEntry, error) {
	response, err := c.apiClient.ClockStartPost(api.ClockStartPostRequest{
		Project: project,
	})
	if err != nil {
		return api.TimeEntry{}, err
	}

	c.setProject(project)
	c.activeEntry = response.Entry
	c.dispatcher.Dispatch(EVENT_CLOCK_START, c.activeEntry)
	return c.activeEntry, nil
}

func (c *clock) Stop() (api.TimeEntry, error) {
	response, err := c.apiClient.ClockStopPost()
	if err != nil {
		return api.TimeEntry{}, err
	}

	c.activeEntry = api.TimeEntry{}
	c.dispatcher.Dispatch(EVENT_CLOCK_STOP, c.activeEntry)
	return response.Entry, nil
}

func (c *clock) setProject(project string) {
	locstor.SetItem("active_project", project)
}
