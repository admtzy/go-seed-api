package main

import (
	"log"
	"net/http"
	"time"

	"go-seed-api/database"
	"go-seed-api/routes"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("db connect:", err)
	}

	router := routes.RegisterRoutes()

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server running at http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
