package routes

import (
	"net/http"

	"autotransport/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux) {
	// CARS
	r.Get("/api/cars", handlers.GetCars)
	r.Post("/api/cars", handlers.AddCarModel)
	r.Put("/api/cars/{id}", handlers.UpdateCarModelByID)
	r.Delete("/api/cars/{id}", handlers.DeleteCarModel)

	// BOOKINGS
	r.Get("/api/bookings", handlers.GetBookings)
	r.Post("/api/bookings", handlers.AddBooking)
	r.Put("/api/bookings/{id}", handlers.UpdateBookingByID)
	r.Delete("/api/bookings/{id}", handlers.DeleteBooking)

	// WEB PAGES
	r.Get("/web/dashboard", serveDashboard)
	r.Get("/web/trips", serveTrips)
	r.Get("/web/cars", serveCars)
	r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

	r.Get("/web", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusMovedPermanently)
	})
}

func serveDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/dashboard.html")
}

func serveTrips(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/trips.html")
}

func serveCars(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/cars.html")
}
