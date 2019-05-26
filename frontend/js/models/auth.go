package models

import (
	"errors"

	"github.com/go-humble/locstor"
	"github.com/winded/tyomaa/frontend/js/dispatcher"
	"github.com/winded/tyomaa/shared/api"
	"github.com/winded/tyomaa/shared/api/client"
)

var (
	NotLoggedInErr = errors.New("Not logged in")
)

type Auth interface {
	LoadToken()
	GetToken() string
	GetUser() (api.User, error)

	Login(username, password string) error
	Logout()
}

type auth struct {
	apiClient  *client.ApiClient
	dispatcher *dispatcher.Dispatcher

	token string
	user  api.User
}

func NewAuth(apiClient *client.ApiClient, disp *dispatcher.Dispatcher) Auth {
	return &auth{
		apiClient:  apiClient,
		dispatcher: disp,
	}
}

func (a *auth) setToken(token string) {
	a.token = token
	a.apiClient.Settings.Token = token
	locstor.SetItem("token", token)
	a.dispatcher.Dispatch(EVENT_CHANGE_TOKEN, token)
}

func (a *auth) LoadToken() {
	token, _ := locstor.GetItem("token")
	a.setToken(token)
}

func (a *auth) GetToken() string {
	return a.token
}

func (a *auth) GetUser() (api.User, error) {
	if a.user.ID != 0 {
		return a.user, nil
	}

	response, err := a.apiClient.AuthTokenGet()
	if err != nil {
		return api.User{}, err
	}

	if response.User == nil {
		return api.User{}, NotLoggedInErr
	}

	a.user = *response.User
	return a.user, nil
}

func (a *auth) Login(username, password string) error {
	request := api.TokenPostRequest{
		Username: username,
		Password: password,
	}

	response, err := a.apiClient.AuthTokenPost(request)
	if err != nil {
		return err
	}

	a.setToken(response.Token)
	return nil
}

func (a *auth) Logout() {
	a.setToken("")
}
