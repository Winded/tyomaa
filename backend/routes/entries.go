package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/winded/tyomaa/backend/middleware"

	"github.com/jinzhu/gorm"

	"github.com/winded/tyomaa/backend/util"

	"github.com/winded/tyomaa/shared/api"

	"github.com/winded/tyomaa/backend/db"
	"github.com/winded/tyomaa/backend/util/context"

	"github.com/gorilla/mux"
)

func EntriesRoutes(router *mux.Router) {
	router.Use(middleware.Authorization)

	router.Path("").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer util.HandleApiError(w)

		ctx := context.Get(r)

		var entries []db.TimeEntry
		if err := db.Instance.Where("user_id = ?", ctx.User.ID).Order("start DESC").Find(&entries).Error; err != nil {
			panic(err)
		}

		var apiEntries []api.TimeEntry
		for _, entry := range entries {
			apiEntries = append(apiEntries, entry.ToApiFormat())
		}
		if apiEntries == nil {
			apiEntries = []api.TimeEntry{}
		}

		json.NewEncoder(w).Encode(api.EntriesGetResponse{
			Entries: apiEntries,
		})
	})

	router.Path("").Methods(http.MethodPost).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer util.HandleApiError(w)

		ctx := context.Get(r)

		var request api.EntriesPostRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			panic(api.Error(http.StatusBadRequest, "Malformed body"))
		}

		if !util.ValidateNameIdentifier(request.Entry.Project) {
			panic(api.Error(http.StatusBadRequest, "Project: "+util.NameIdentifierError))
		}
		if request.Entry.End != nil && (request.Entry.Start == *request.Entry.End || request.Entry.Start.After(*request.Entry.End)) {
			panic(api.Error(http.StatusBadRequest, "End time must be later than start time"))
		}

		entry := db.TimeEntry{
			UserID:  ctx.User.ID,
			Project: request.Entry.Project,
			Start:   request.Entry.Start,
			End:     request.Entry.End,
		}
		if err := db.Instance.Save(&entry).Error; err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(api.EntriesPostResponse{
			Entry: entry.ToApiFormat(),
		})
	})

	router.Path("/{entryId:[0-9]+}").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer util.HandleApiError(w)

		ctx := context.Get(r)

		vars := mux.Vars(r)
		entryId, err := strconv.Atoi(vars["entryId"])
		if err != nil {
			panic(err)
		}

		var entry db.TimeEntry
		err = db.Instance.Where("id = ? AND user_id = ?", entryId, ctx.User.ID).First(&entry).Error
		if gorm.IsRecordNotFoundError(err) {
			panic(api.Error(http.StatusNotFound, "Entry not found"))
		} else if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(api.EntriesSingleGetResponse{
			Entry: entry.ToApiFormat(),
		})
	})

	router.Path("/{entryId:[0-9]+}").Methods(http.MethodPost).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer util.HandleApiError(w)

		ctx := context.Get(r)

		vars := mux.Vars(r)
		entryId, err := strconv.Atoi(vars["entryId"])
		if err != nil {
			panic(err)
		}

		var request api.EntriesSinglePostRequest
		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			panic(api.Error(http.StatusBadRequest, "Malformed body"))
		}

		if !util.ValidateNameIdentifier(request.Entry.Project) {
			panic(api.Error(http.StatusBadRequest, "Project: "+util.NameIdentifierError))
		}
		if request.Entry.End != nil && (request.Entry.Start == *request.Entry.End || request.Entry.Start.After(*request.Entry.End)) {
			panic(api.Error(http.StatusBadRequest, "End time must be later than start time"))
		}

		var entry db.TimeEntry
		err = db.Instance.Where("id = ? AND user_id = ?", entryId, ctx.User.ID).First(&entry).Error
		if gorm.IsRecordNotFoundError(err) {
			panic(api.Error(http.StatusNotFound, "Entry not found"))
		} else if err != nil {
			panic(err)
		}

		entry.Project = request.Entry.Project
		entry.Start = request.Entry.Start
		entry.End = request.Entry.End
		if err := db.Instance.Save(&entry); err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(api.EntriesSinglePostResponse{
			Entry: entry.ToApiFormat(),
		})
	})
}
