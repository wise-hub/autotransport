package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"autotransport/internal/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
)

func verboseErrorLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		if ww.Status() >= 400 {
			log.Printf("Error: %s %s - Status: %d, User-Agent: %s, Remote IP: %s\n",
				r.Method, r.RequestURI, ww.Status(), r.UserAgent(), r.RemoteAddr)
		}
	})
}

func initDB() *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"postgres")

	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to the database: %v", err))
	}

	var exists bool
	err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", "autotransport").Scan(&exists)
	if err != nil {
		panic(fmt.Sprintf("Unable to check if database exists: %v", err))
	}

	if !exists {
		_, err = db.Exec(context.Background(), "CREATE DATABASE autotransport")
		if err != nil {
			panic(fmt.Sprintf("Unable to create database: %v", err))
		}
	}

	db.Close()

	dsn = fmt.Sprintf("postgres://%s:%s@localhost:5432/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"autotransport")

	db, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to the database: %v", err))
	}

	sqlFile, err := os.ReadFile("./ddl.sql")
	if err != nil {
		panic(fmt.Sprintf("Unable to read SQL file: %v", err))
	}

	_, err = db.Exec(context.Background(), string(sqlFile))
	if err != nil {
		panic(fmt.Sprintf("Failed to execute DDL statements: %v", err))
	}

	return db
}

func main() {
	db := initDB()
	defer db.Close()

	api.SetDB(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(verboseErrorLogger)

	// CARS
	r.Get("/api/cars", api.GetCars)
	r.Post("/api/cars", api.AddCarModel)
	r.Put("/api/cars/{id}", api.UpdateCarModelByID)
	r.Delete("/api/cars/{id}", api.DeleteCarModel)

	// BOOKINGS
	r.Get("/api/bookings", api.GetBookings)
	r.Post("/api/bookings", api.AddBooking)
	r.Put("/api/bookings/{id}", api.UpdateBookingByID)
	r.Delete("/api/bookings/{id}", api.DeleteBooking)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	if err := http.ListenAndServe(":8085", r); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
