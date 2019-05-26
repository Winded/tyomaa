package app

import (
	"github.com/winded/tyomaa/frontend/js/dispatcher"
	"github.com/winded/tyomaa/frontend/js/models"
	"github.com/winded/tyomaa/shared/api/client"
)

type Context struct {
	Dispatcher *dispatcher.Dispatcher

	Auth    models.Auth
	Clock   models.Clock
	Entries models.Entries
}

func NewContext(disp *dispatcher.Dispatcher, apiClient *client.ApiClient) *Context {
	return &Context{
		Dispatcher: disp,
		Auth:       models.NewAuth(apiClient, disp),
		Clock:      models.NewClock(apiClient, disp),
		Entries:    models.NewEntries(apiClient, disp),
	}
}
