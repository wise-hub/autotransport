package main

import (
	"fmt"
	"net/http"

	"autotransport/internal/db"
	"autotransport/internal/handlers"
	"autotransport/internal/middleware"
	"autotransport/internal/routes"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	handlers.SetDB(database)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.VerboseErrorLogger)

	routes.RegisterRoutes(r)

	if err := http.ListenAndServe(":8085", r); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
