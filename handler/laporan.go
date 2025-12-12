package handler

import (
	"go-seed-api/database"
	"go-seed-api/middleware"
	"go-seed-api/models"
	"go-seed-api/utils"
	"net/http"
)

// Pure function: transformasi StokHistory + nama bibit + user ke JSON
func mapToJSON(sh models.StokHistory, bibitNama, userNama string) map[string]interface{} {
	return map[string]interface{}{
		"id":         sh.ID,
		"bibit_id":   sh.BibitID,
		"bibit_nama": bibitNama,
		"user_id":    sh.UserID,
		"user_nama":  userNama,
		"tipe":       sh.Tipe,
		"jumlah":     sh.Jumlah,
		"created_at": sh.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func GetLaporan(w http.ResponseWriter, r *http.Request) {
	// Ambil user dari context
	claims, _ := r.Context().Value(middleware.UserKey).(map[string]interface{})
	username := ""
	if claims != nil {
		username = claims["username"].(string)
	}

	// Ambil data dari DB
	rows, err := database.DB.Query(`
		SELECT 
			s.id,
			s.bibit_id,
			b.nama AS bibit_nama,
			s.user_id,
			u.username AS user_nama,
			s.tipe,
			s.jumlah,
			s.created_at
		FROM stok_history s
		JOIN bibit b ON b.id = s.bibit_id
		JOIN users u ON u.id = s.user_id
		ORDER BY s.created_at DESC
	`)
	if err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	// Slice untuk menampung semua row sementara
	type rowData struct {
		SH        models.StokHistory
		BibitNama string
		UserNama  string
	}
	var rowsData []rowData

	for rows.Next() {
		var sh models.StokHistory
		var bibitNama, userNama string

		if err := rows.Scan(
			&sh.ID,
			&sh.BibitID,
			&bibitNama,
			&sh.UserID,
			&userNama,
			&sh.Tipe,
			&sh.Jumlah,
			&sh.CreatedAt,
		); err != nil {
			utils.WriteError(w, 500, err.Error())
			return
		}

		rowsData = append(rowsData, rowData{sh, bibitNama, userNama})
	}

	// Transformasi → JSON (pure + map)
	result := utils.Map(rowsData, func(r rowData) map[string]interface{} {
		return mapToJSON(r.SH, r.BibitNama, r.UserNama)
	})

	// Contoh reduce: total jumlah tipe "masuk"
	totalMasuk := utils.Reduce(result, 0, func(acc int, r map[string]interface{}) int {
		if r["tipe"] == "masuk" {
			return acc + r["jumlah"].(int)
		}
		return acc
	})

	utils.WriteJSON(w, 200, utils.JSON{
		"laporan":     result,
		"total_masuk": totalMasuk,
		"user":        username,
	})
}
