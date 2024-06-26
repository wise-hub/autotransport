package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CarBooking struct {
	CarModel string    `json:"car_model"`
	Bookings []Booking `json:"bookings"`
}

type Booking struct {
	ID           int    `json:"id"`
	CarModelID   int    `json:"car_model_id"`
	CarModelName string `json:"car_model_name"`
	Department   string `json:"department"`
	UsageDate    string `json:"usage_date"`
	Destination  string `json:"destination"`
}

type BookingRequest struct {
	CarModelID  int    `json:"car_model_id"`
	Department  string `json:"department"`
	UsageDate   string `json:"usage_date"`
	Destination string `json:"destination"`
}

type CarModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *pgxpool.Pool

func SetDB(database *pgxpool.Pool) {
	db = database
}

func GetBookings(w http.ResponseWriter, r *http.Request) {
	query := `SELECT cm.name, b.id, b.car_model_id, b.department, b.usage_date, b.destination
		FROM bookings b
		JOIN car_models cm ON cm.id = b.car_model_id 
		WHERE b.deleted = FALSE 
		ORDER BY cm.name, b.usage_date`
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Unable to retrieve bookings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	carBookingsMap := make(map[string]*CarBooking)

	for rows.Next() {
		var carModel string
		var booking Booking

		if err := rows.Scan(&carModel, &booking.ID, &booking.CarModelID, &booking.Department, &booking.UsageDate, &booking.Destination); err != nil {
			http.Error(w, "Unable to scan data", http.StatusInternalServerError)
			return
		}

		booking.CarModelName = carModel

		if carBooking, exists := carBookingsMap[carModel]; exists {
			carBooking.Bookings = append(carBooking.Bookings, booking)
		} else {
			carBookingsMap[carModel] = &CarBooking{
				CarModel: carModel,
				Bookings: []Booking{booking},
			}
		}
	}

	var carBookings []CarBooking
	for _, carBooking := range carBookingsMap {
		carBookings = append(carBookings, *carBooking)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(carBookings)
}

func AddBooking(w http.ResponseWriter, r *http.Request) {
	var newBooking Booking
	err := json.NewDecoder(r.Body).Decode(&newBooking)
	if err != nil {
		http.Error(w, "Invalid input - "+err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO bookings (car_model_id, department, usage_date, destination) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(context.Background(), query, newBooking.CarModelID, newBooking.Department, newBooking.UsageDate, newBooking.Destination)
	if err != nil {
		http.Error(w, "Unable to add booking", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func AddCarModel(w http.ResponseWriter, r *http.Request) {
	var newCarModel CarModel
	err := json.NewDecoder(r.Body).Decode(&newCarModel)
	if err != nil {
		http.Error(w, "Invalid input - "+err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO car_models (name) VALUES ($1)`
	_, err = db.Exec(context.Background(), query, newCarModel.Name)
	if err != nil {
		http.Error(w, "Unable to add car model", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func UpdateCarModelByID(w http.ResponseWriter, r *http.Request) {
	carModelID := chi.URLParam(r, "id")
	var updatedCarModel CarModel
	err := json.NewDecoder(r.Body).Decode(&updatedCarModel)
	if err != nil {
		http.Error(w, "Invalid input - "+err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE car_models SET name = $1 WHERE id = $2 AND deleted = FALSE`
	_, err = db.Exec(context.Background(), query, updatedCarModel.Name, carModelID)
	if err != nil {
		http.Error(w, "Unable to update car model", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateBookingByID(w http.ResponseWriter, r *http.Request) {
	bookingID := chi.URLParam(r, "id")
	var updatedBooking BookingRequest
	err := json.NewDecoder(r.Body).Decode(&updatedBooking)
	if err != nil {
		http.Error(w, "Invalid input - "+err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE bookings SET car_model_id = $1, department = $2, usage_date = $3, destination = $4 WHERE id = $5 AND deleted = FALSE`
	_, err = db.Exec(context.Background(), query, updatedBooking.CarModelID, updatedBooking.Department, updatedBooking.UsageDate, updatedBooking.Destination, bookingID)
	if err != nil {
		http.Error(w, "Unable to update booking", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetCars(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name FROM car_models where deleted = false`
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Unable to retrieve cars", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cars []CarModel

	for rows.Next() {
		var car CarModel

		if err := rows.Scan(&car.ID, &car.Name); err != nil {
			http.Error(w, "Unable to scan data", http.StatusInternalServerError)
			return
		}

		cars = append(cars, car)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func DeleteCarModel(w http.ResponseWriter, r *http.Request) {
	carModelID := chi.URLParam(r, "id")
	query := `UPDATE car_models SET deleted = TRUE WHERE id = $1`
	_, err := db.Exec(context.Background(), query, carModelID)
	if err != nil {
		http.Error(w, "Unable to mark car model as deleted", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteBooking(w http.ResponseWriter, r *http.Request) {
	bookingID := chi.URLParam(r, "id")
	query := `UPDATE bookings SET deleted = TRUE WHERE id = $1`
	_, err := db.Exec(context.Background(), query, bookingID)
	if err != nil {
		http.Error(w, "Unable to mark booking as deleted", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
