package context

import (
	"context"
	"net/http"

	"github.com/winded/tyomaa/backend/db"
	"github.com/winded/tyomaa/backend/util/token"
)

const (
	contextDataKey = "contextData"
)

type ContextData struct {
	Token  string
	Claims token.Claims
	User   db.User
}

func Get(r *http.Request) ContextData {
	data, ok := r.Context().Value(contextDataKey).(ContextData)
	if ok {
		return data
	} else {
		return ContextData{}
	}
}

func Set(r *http.Request, data ContextData) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contextDataKey, data))
}
