package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/winded/tyomaa/backend/db"
	"github.com/winded/tyomaa/backend/middleware"
	"github.com/winded/tyomaa/backend/util"
	"github.com/winded/tyomaa/backend/util/context"
	"github.com/winded/tyomaa/shared/api"
)

func ClockRoutes(router *mux.Router) {
	router.Use(middleware.Authorization)

	router.Path("").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer util.HandleApiError(w)

		ctx := context.Get(r)

		var activeEntry db.TimeEntry
		err := db.Instance.First(&activeEntry, `user_id = ? AND "end" IS NULL`, ctx.User.ID).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			panic(err)
		}

		var response api.ClockGetResponse
		if activeEntry.ID != 0 {
			apiEntry := activeEntry.ToApiFormat()
			response.Entry = &apiEntry
		}

		json.NewEncoder(w).Encode(response)
	})

	router.Path("/start").Methods(http.MethodPost).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer util.HandleApiError(w)

		ctx := context.Get(r)

		var request api.ClockStartPostRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			panic(api.Error(http.StatusBadRequest, "Malformed body"))
		}

		if !util.ValidateNameIdentifier(request.Project) {
			panic(api.Error(http.StatusBadRequest, "Project: "+util.NameIdentifierError))
		}

		activeEntry := db.TimeEntry{
			UserID:  ctx.User.ID,
			Project: request.Project,
			Start:   time.Now(),
			End:     nil,
		}
		if err := db.Instance.Save(&activeEntry).Error; err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(api.ClockStartPostResponse{
			Entry: activeEntry.ToApiFormat(),
		})
	})

	router.Path("/stop").Methods(http.MethodPost).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer util.HandleApiError(w)

		ctx := context.Get(r)

		var activeEntry db.TimeEntry
		err := db.Instance.Where(`user_id = ? AND "end" IS NULL`, ctx.User.ID).First(&activeEntry).Error
		if gorm.IsRecordNotFoundError(err) {
			panic(api.Error(http.StatusBadRequest, "Active entry does not exist"))
		} else if err != nil {
			panic(err)
		}

		end := time.Now()
		activeEntry.End = &end
		if err := db.Instance.Save(&activeEntry).Error; err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(api.ClockStopPostResponse{
			Entry: activeEntry.ToApiFormat(),
		})
	})
}
