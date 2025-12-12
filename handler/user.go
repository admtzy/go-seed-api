package handler

import (
	"encoding/json"
	"go-seed-api/database"
	"go-seed-api/models"
	"go-seed-api/utils"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	_, err := database.DB.Exec(`
        INSERT INTO users (username, password, role)
        VALUES ($1, $2, 'user')
    `, u.Username, u.Password)
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

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	var user models.User
	err := database.DB.QueryRow(`
        SELECT id, username, password, role
        FROM users
        WHERE username = $1
    `, loginData.Username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		http.Error(w, "user not found", 404)
		return
	}

	if loginData.Password != user.Password {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, "gagal generate token", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "login success",
		"token":   token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
        SELECT id, username, role, created_at
        FROM users
        ORDER BY id ASC
    `)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}
