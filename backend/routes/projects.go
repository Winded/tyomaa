package routes

import (
	"encoding/json"
	"net/http"

	"github.com/peterhellberg/duration"

	"github.com/winded/tyomaa/backend/db"
	"github.com/winded/tyomaa/shared/api"

	"github.com/winded/tyomaa/backend/util"
	"github.com/winded/tyomaa/backend/util/context"

	"github.com/gorilla/mux"
	"github.com/winded/tyomaa/backend/middleware"
)

func ProjectsRoutes(router *mux.Router) {
	const projectsQuery = `SELECT "project" as "name", SUM("end" - "start") as "total_time" FROM time_entries WHERE user_id = ? GROUP BY "project" ORDER BY "project" ASC`

	router.Use(middleware.Authorization)

	router.Path("").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer util.HandleApiError(w)

		ctx := context.Get(r)

		rows, err := db.Instance.Raw(projectsQuery, ctx.User.ID).Rows()
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var results []api.Project
		for rows.Next() {
			var (
				name string
				t    string
			)
			if err := rows.Scan(&name, &t); err != nil {
				panic(err)
			}

			totalTime, err := duration.Parse(t)
			if err != nil {
				panic(err)
			}

			results = append(results, api.Project{
				Name:      name,
				TotalTime: totalTime,
			})
		}

		json.NewEncoder(w).Encode(api.ProjectsGetResponse{
			Projects: results,
		})
	})
}
