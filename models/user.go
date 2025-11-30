package models

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
