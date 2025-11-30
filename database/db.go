package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect melakukan side-effect: membuka koneksi DB
func Connect() error {
	connStr := "host=localhost user=postgres password=1234 dbname=bibitdb sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	log.Println("Database connected")
	return nil
}
