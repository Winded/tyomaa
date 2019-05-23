package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/winded/tyomaa/backend/db"
	"github.com/winded/tyomaa/backend/middleware"
	"github.com/winded/tyomaa/backend/routes"
	"github.com/winded/tyomaa/backend/util"
)

func main() {
	dbConn, err := db.Init()
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	db.AutoMigrate(dbConn)

	app := mux.NewRouter()
	app.Use(middleware.Json)
	app.Use(middleware.TokenSession)

	baseRouter := app.PathPrefix(util.EnvOrDefault("ROOT_URL", "")).Subrouter()
	baseRouter.Use(middleware.Authentication)

	routes.AuthRoutes(baseRouter.PathPrefix("/auth").Subrouter())
	routes.EntriesRoutes(baseRouter.PathPrefix("/entries").Subrouter())
	routes.ClockRoutes(baseRouter.PathPrefix("/clock").Subrouter())
	routes.ProjectsRoutes(baseRouter.Path("/projects").Subrouter())

	fmt.Println("Starting web server...")
	http.ListenAndServe(":80", app)
}
