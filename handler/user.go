package handler

import (
	"encoding/json"
	"go-seed-api/database"
	"go-seed-api/models"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)

	// 👇 Tidak ada hashing
	_, err := database.DB.Exec(
		`INSERT INTO users (username,password,role) VALUES ($1,$2,$3)`,
		u.Username, u.Password, "user",
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "user registered",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&loginData)

	var user models.User

	err := database.DB.QueryRow(
		"SELECT id, password, role FROM users WHERE username=$1",
		loginData.Username,
	).Scan(&user.ID, &user.Password, &user.Role)

	if err != nil {
		http.Error(w, "user not found", 404)
		return
	}

	// 👇 Perbandingan langsung, tidak pakai bcrypt
	if loginData.Password != user.Password {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "login success",
		"user_id": user.ID,
		"role":    user.Role,
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, username, role FROM users")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var users []models.User

	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Username, &u.Role)
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}
