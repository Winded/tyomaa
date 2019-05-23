package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/winded/tyomaa/backend/services"
	"github.com/winded/tyomaa/backend/util"
	"github.com/winded/tyomaa/backend/util/context"
	"github.com/winded/tyomaa/shared/api"
)

func AuthRoutes(router *mux.Router) {
	router.Path("/token").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response api.TokenGetResponse

		ctx := context.Get(r)
		response.Token = ctx.Token
		if ctx.User.ID != 0 {
			user := ctx.User.ToApiFormat()
			response.User = &user
		}

		json.NewEncoder(w).Encode(&response)
	})

	router.Path("/token").Methods(http.MethodPost).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer util.HandleApiError(w)

		var request api.TokenPostRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			panic(api.Error(http.StatusBadRequest, "Malformed body"))
		}

		user := services.AuthenticateUser(request.Username, request.Password)
		if user == nil {
			panic(api.Error(http.StatusBadRequest, "Invalid username or password"))
		}

		tk, err := services.CreateToken(user)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(api.TokenPostResponse{
			Token: tk,
		})
	})
}
