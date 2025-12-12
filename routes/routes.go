package routes

import (
	"go-seed-api/handler"
	"go-seed-api/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// Global middleware
	r.Use(middleware.CORS)

	r.Methods(http.MethodOptions).HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	// PUBLIC ROUTES
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")

	// PROTECTED ROUTES (WAJIB LOGIN)
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.Auth)

	// User
	protected.HandleFunc("/users", handler.GetUsers).Methods("GET")

	// Bibit
	protected.HandleFunc("/bibit", handler.GetBibit).Methods("GET")
	protected.HandleFunc("/bibit", handler.CreateBibit).Methods("POST")

	// Stok
	protected.HandleFunc("/stok/{id:[0-9]+}", handler.UpdateStok).Methods("PUT")

	// Rekomendasi
	protected.HandleFunc("/rekomendasi", handler.GetRekomendasi).Methods("GET")

	// Laporan
	protected.HandleFunc("/laporan", handler.GetLaporan).Methods("GET")

	return r
}
