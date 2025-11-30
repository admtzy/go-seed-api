package models

type RekomendasiLog struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Tanah     string `json:"tanah"`
	Curah     int    `json:"curah"`
	Luas      int    `json:"luas"`
	Hasil     string `json:"hasil"`
	CreatedAt string `json:"created_at"`
}
