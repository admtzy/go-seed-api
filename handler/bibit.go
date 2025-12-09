package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-seed-api/database"
	"go-seed-api/models"
	"go-seed-api/utils"

	"github.com/gorilla/mux"
)

/* ===========================================================
   F U N C T I O N A L   P R O G R A M M I N G   S E C T I O N
=========================================================== */

// Pure function untuk hitung kebutuhan bibit (rekursi)
func HitungKebutuhanBibit(luas float64) int {
	if luas <= 0 {
		return 0
	}
	return 5 + HitungKebutuhanBibit(luas-0.1)
}

// HOF Validator Generic
type Validator func() (bool, string)

func ValidateNotEmpty(field, name string) Validator {
	return func() (bool, string) {
		if field == "" {
			return false, fmt.Sprintf("%s tidak boleh kosong", name)
		}
		return true, ""
	}
}

func ValidatePositive(value int, name string) Validator {
	return func() (bool, string) {
		if value < 0 {
			return false, fmt.Sprintf("%s harus positif", name)
		}
		return true, ""
	}
}

/* ===========================================================
	API CONTROLLER
=========================================================== */

// CREATE BIBIT (Pure logic terpisah dari side-effect)
func CreateBibit(w http.ResponseWriter, r *http.Request) {
	var bibit models.Bibit

	if err := json.NewDecoder(r.Body).Decode(&bibit); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Format JSON salah: "+err.Error())
		return
	}

	// HOF Validasi
	validations := []Validator{
		ValidateNotEmpty(bibit.Nama, "Nama bibit"),
		ValidateNotEmpty(bibit.Kualitas, "Kualitas"),
		ValidatePositive(bibit.Stok, "Stok"),
	}

	for _, v := range validations {
		if ok, msg := v(); !ok {
			utils.WriteError(w, http.StatusBadRequest, msg)
			return
		}
	}

	now := time.Now()
	const q = `
		INSERT INTO bibit (nama, kualitas, stok, tanah, curah_hujan, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$6)`

	_, err := database.DB.Exec(q,
		bibit.Nama, bibit.Kualitas, bibit.Stok,
		bibit.Tanah, bibit.CurahHujan, now)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.JSON{"message": "Bibit berhasil ditambahkan"})
}

// GET BIBIT (Map + Filter + Reduce FP COMPLETE)
func GetBibit(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, nama, kualitas, stok, tanah, curah_hujan, created_at, updated_at 
		FROM bibit`)
	if err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	var list []models.Bibit
	for rows.Next() {
		var b models.Bibit
		rows.Scan(&b.ID, &b.Nama, &b.Kualitas, &b.Stok,
			&b.Tanah, &b.CurahHujan, &b.CreatedAt, &b.UpdatedAt)
		list = append(list, b.Clone())
	}

	// Filter stok > 0
	aktif := utils.Filter(list, func(b models.Bibit) bool {
		return b.Stok > 0
	})

	// Map → Summary ONLY
	type Summary struct {
		ID   int    `json:"id"`
		Nama string `json:"nama"`
		Stok int    `json:"stok"`
	}

	result := utils.Map(aktif, func(b models.Bibit) Summary {
		return Summary{b.ID, b.Nama, b.Stok}
	})

	// Reduce → Total Stok
	total := utils.Reduce(aktif, 0, func(acc int, s models.Bibit) int {
		return acc + s.Stok
	})

	utils.WriteJSON(w, 200, utils.JSON{
		"data":       result,
		"total_stok": total,
	})
}

// UPDATE STOK (Side-effect terkendali di 1 tempat)
func UpdateStok(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var body struct {
		Delta int `json:"delta"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	var current int
	err := tx.QueryRow(`
		SELECT stok FROM bibit WHERE id=$1 FOR UPDATE`, id).Scan(&current)

	if err == sql.ErrNoRows {
		utils.WriteError(w, 404, "Bibit tidak ditemukan")
		return
	}

	newStock := current + body.Delta
	if newStock < 0 {
		utils.WriteError(w, 400, "Stok tidak cukup")
		return
	}

	// Update stok
	tx.Exec(`UPDATE bibit SET stok=$1,updated_at=$2 WHERE id=$3`,
		newStock, time.Now(), id)

	// Log history stok
	tipe := utils.If(body.Delta < 0, "keluar", "masuk")
	tx.Exec(`INSERT INTO stok_history (bibit_id, tipe, jumlah) VALUES ($1,$2,$3)`,
		id, tipe, body.Delta)

	tx.Commit()

	utils.WriteJSON(w, 200, utils.JSON{
		"message": "Stok diperbarui",
		"stok":    newStock,
	})
}

// GET REKOMENDASI (Closure + Pure + Recursion)
func GetRekomendasi(w http.ResponseWriter, r *http.Request) {
	soil := r.URL.Query().Get("tanah")
	rain, _ := strconv.Atoi(r.URL.Query().Get("curah"))
	luas, _ := strconv.ParseFloat(r.URL.Query().Get("luas"), 64)

	// Closure ID request
	counter := func() func() string {
		count := 0
		return func() string {
			count++
			return fmt.Sprintf("REQ-%d", count)
		}
	}()
	requestID := counter()

	kebutuhan := HitungKebutuhanBibit(luas)

	utils.WriteJSON(w, 200, utils.JSON{
		"id":          requestID,
		"tanah":       soil,
		"curah_hujan": rain,
		"kebutuhan":   kebutuhan,
	})
}
