package handler

import (
	"database/sql"
	"go-seed-api/database"
	"go-seed-api/utils"
	"net/http"
	"time"
)

// LaporanEntry representasi immutable untuk laporan stok
type LaporanEntry struct {
	Bibit  string    `json:"bibit"`
	Tipe   string    `json:"tipe"`
	Jumlah int       `json:"jumlah"`
	Waktu  time.Time `json:"waktu"`
}

// Pure function: baca rows dan ubah menjadi slice LaporanEntry
func parseLaporanRows(rows *sql.Rows) ([]LaporanEntry, error) {
	var result []LaporanEntry

	for rows.Next() {
		var nama, tipe string
		var jumlah int
		var created time.Time

		if err := rows.Scan(&nama, &tipe, &jumlah, &created); err != nil {
			return nil, err
		}

		entry := LaporanEntry{
			Bibit:  nama,
			Tipe:   tipe,
			Jumlah: jumlah,
			Waktu:  created,
		}

		// Immutable append (FP style)
		result = append(result, entry)
	}

	return result, nil
}

// GET /laporan
func GetLaporan(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT b.nama, s.tipe, s.jumlah, s.created_at
		FROM stok_history s
		JOIN bibit b ON b.id = s.bibit_id
		ORDER BY s.created_at DESC`)
	if err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	// Pure function transformasi rows
	laporan, err := parseLaporanRows(rows)
	if err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}

	// Contoh Filter: hanya "masuk"
	filtered := utils.Filter(laporan, func(e LaporanEntry) bool {
		return e.Tipe == "masuk"
	})

	// Contoh Map: hanya ambil nama dan jumlah
	mapped := utils.Map(filtered, func(e LaporanEntry) map[string]interface{} {
		return map[string]interface{}{
			"bibit":  e.Bibit,
			"jumlah": e.Jumlah,
		}
	})

	// Contoh Reduce: total jumlah
	total := utils.Reduce(filtered, 0, func(acc int, e LaporanEntry) int {
		return acc + e.Jumlah
	})

	utils.WriteJSON(w, 200, utils.JSON{
		"data":        mapped,
		"total_masuk": total,
	})
}
