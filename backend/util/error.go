package util

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/winded/tyomaa/shared/api"
)

func WriteApiError(w http.ResponseWriter, err error) {
	apiErr, ok := err.(api.ApiError)
	if ok {
		w.WriteHeader(apiErr.Status)
		json.NewEncoder(w).Encode(apiErr)
	} else {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.Error(http.StatusInternalServerError,
			"An internal server error occurred"))
	}
}

func HandleApiError(w http.ResponseWriter) {
	rec := recover()
	if rec == nil {
		return
	}

	err := rec.(error)
	WriteApiError(w, err)
}
