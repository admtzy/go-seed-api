package routes

import (
	"go-seed-api/handler" // ⬅️ hanya ini yang diganti

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// --- AUTH / USER MANAGEMENT ---
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/users", handler.GetUsers).Methods("GET")

	// --- BIBIT CRUD ---
	r.HandleFunc("/bibit", handler.CreateBibit).Methods("POST")
	r.HandleFunc("/bibit", handler.GetBibit).Methods("GET")
	// r.HandleFunc("/bibit/{id:[0-9]+}", handler.DeleteBibit).Methods("DELETE")

	// --- UPDATE STOK & HISTORY ---
	r.HandleFunc("/stok/{id:[0-9]+}", handler.UpdateStok).Methods("PUT")

	// --- REKOMENDASI ---
	r.HandleFunc("/rekomendasi", handler.GetRekomendasi).Methods("GET")

	// --- LAPORAN ---
	r.HandleFunc("/laporan", handler.GetLaporan).Methods("GET")

	return r
}
