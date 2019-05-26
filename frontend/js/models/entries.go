package models

import (
	"github.com/winded/tyomaa/frontend/js/dispatcher"
	"github.com/winded/tyomaa/shared/api"
	"github.com/winded/tyomaa/shared/api/client"
)

type Entries interface {
	List() ([]api.TimeEntry, error)
	Get(entryId uint) (api.TimeEntry, error)

	Create(entry api.TimeEntry) (api.TimeEntry, error)
	Update(updatedEntry api.TimeEntry) (api.TimeEntry, error)
	Delete(entryId uint) (api.TimeEntry, error)
}

type entries struct {
	apiClient  *client.ApiClient
	dispatcher *dispatcher.Dispatcher
}

func NewEntries(apiClient *client.ApiClient, disp *dispatcher.Dispatcher) Entries {
	return &entries{
		apiClient:  apiClient,
		dispatcher: disp,
	}
}

func (e *entries) List() ([]api.TimeEntry, error) {
	response, err := e.apiClient.EntriesGet(api.EntriesGetRequest{})
	if err != nil {
		return nil, err
	}

	return response.Entries, nil
}

func (e *entries) Get(entryId uint) (api.TimeEntry, error) {
	response, err := e.apiClient.EntriesSingleGet(entryId)
	if err != nil {
		return api.TimeEntry{}, err
	}

	return response.Entry, nil
}

func (e *entries) Create(entry api.TimeEntry) (api.TimeEntry, error) {
	response, err := e.apiClient.EntriesPost(api.EntriesPostRequest{
		Entry: entry,
	})
	if err != nil {
		return api.TimeEntry{}, err
	}

	e.dispatcher.Dispatch(EVENT_ENTRY_CREATED, response.Entry)
	return response.Entry, nil
}

func (e *entries) Update(updatedEntry api.TimeEntry) (api.TimeEntry, error) {
	response, err := e.apiClient.EntriesSinglePost(updatedEntry.ID, api.EntriesSinglePostRequest{
		Entry: updatedEntry,
	})
	if err != nil {
		return api.TimeEntry{}, err
	}

	e.dispatcher.Dispatch(EVENT_ENTRY_UPDATED, response.Entry)
	return response.Entry, nil
}

func (e *entries) Delete(entryId uint) (api.TimeEntry, error) {
	response, err := e.apiClient.EntriesSingleDelete(entryId)
	if err != nil {
		return api.TimeEntry{}, err
	}

	e.dispatcher.Dispatch(EVENT_ENTRY_DELETED, response.Entry)
	return response.Entry, nil
}
