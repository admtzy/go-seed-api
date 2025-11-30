package models

type StokHistory struct {
	ID        int    `json:"id"`
	BibitID   int    `json:"bibit_id"`
	UserID    int    `json:"user_id"`
	Tipe      string `json:"tipe"`
	Jumlah    int    `json:"jumlah"`
	CreatedAt string `json:"created_at"`
}
