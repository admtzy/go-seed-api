package models

import "time"

type Bibit struct {
	ID         int       `json:"id"`
	Nama       string    `json:"nama"`
	Kualitas   string    `json:"kualitas"`
	Stok       int       `json:"stok"`
	Tanah      string    `json:"tanah"`
	CurahHujan int       `json:"curah_hujan"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Clone mengembalikan salinan (praktik immutable-ish)
func (b Bibit) Clone() Bibit {
	return Bibit{
		ID:         b.ID,
		Nama:       b.Nama,
		Kualitas:   b.Kualitas,
		Stok:       b.Stok,
		Tanah:      b.Tanah,
		CurahHujan: b.CurahHujan,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
	}
}
