package handler

import (
	"database/sql"
	"encoding/json"
	"go-seed-api/database"
	"go-seed-api/middleware"
	"go-seed-api/models"
	"go-seed-api/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// BIBIT HANDLER
func CreateBibit(w http.ResponseWriter, r *http.Request) {
	var bibit models.Bibit
	if err := json.NewDecoder(r.Body).Decode(&bibit); err != nil {
		utils.WriteError(w, 400, "invalid body")
		return
	}

	validations := []utils.Validator{
		utils.ValidateNotEmpty(bibit.Nama, "Nama bibit"),
		utils.ValidateNotEmpty(bibit.Kualitas, "Kualitas"),
		utils.ValidateNotEmpty(bibit.Tanah, "Jenis tanah"),
		utils.ValidatePositive(bibit.Stok, "Stok"),
		utils.ValidatePositive(bibit.CurahHujan, "Curah hujan"),
	}

	for _, v := range validations {
		if ok, msg := v(); !ok {
			utils.WriteError(w, 400, msg)
			return
		}
	}

	now := time.Now()
	const q = `
		INSERT INTO bibit (nama,kualitas,stok,tanah,curah_hujan,created_at,updated_at)
		VALUES($1,$2,$3,$4,$5,$6,$6)
	`
	if _, err := database.DB.Exec(q, bibit.Nama, bibit.Kualitas, bibit.Stok, bibit.Tanah, bibit.CurahHujan, now); err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}

	utils.WriteJSON(w, 201, map[string]string{"message": "Bibit berhasil ditambahkan"})
}

func GetBibit(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`SELECT id,nama,kualitas,stok,tanah,curah_hujan,created_at,updated_at FROM bibit`)
	if err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	var list []models.Bibit
	for rows.Next() {
		var b models.Bibit
		rows.Scan(&b.ID, &b.Nama, &b.Kualitas, &b.Stok, &b.Tanah, &b.CurahHujan, &b.CreatedAt, &b.UpdatedAt)
		list = append(list, b.Clone())
	}

	aktif := utils.Filter(list, func(b models.Bibit) bool { return b.Stok > 0 })
	result := utils.Map(aktif, utils.ToSummary)
	total := utils.Reduce(aktif, 0, func(acc int, s models.Bibit) int { return acc + s.Stok })

	utils.WriteJSON(w, 200, map[string]interface{}{"data": result, "total_stok": total})
}

func UpdateStok(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var body struct {
		Delta int `json:"delta"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	claims, ok := r.Context().Value(middleware.UserKey).(map[string]interface{})
	if !ok {
		utils.WriteError(w, 500, "Gagal membaca user dari context")
		return
	}
	userID := int(claims["user_id"].(float64))

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	var current int
	if err := tx.QueryRow(`SELECT stok FROM bibit WHERE id=$1 FOR UPDATE`, id).Scan(&current); err == sql.ErrNoRows {
		utils.WriteError(w, 404, "Bibit tidak ditemukan")
		return
	}

	newStock := current + body.Delta
	if newStock < 0 {
		utils.WriteError(w, 400, "Stok tidak cukup")
		return
	}

	_, err := tx.Exec(`UPDATE bibit SET stok=$1, updated_at=$2 WHERE id=$3`, newStock, time.Now(), id)
	if err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}

	tipe := utils.If(body.Delta < 0, "keluar", "masuk")
	_, err = tx.Exec(`INSERT INTO stok_history (bibit_id,user_id,tipe,jumlah) VALUES($1,$2,$3,$4)`, id, userID, tipe, body.Delta)
	if err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}

	tx.Commit()
	utils.WriteJSON(w, 200, map[string]interface{}{"message": "Stok diperbarui", "stok": newStock})
}

func GetRekomendasi(w http.ResponseWriter, r *http.Request) {
	soil := r.URL.Query().Get("tanah")
	rain, _ := strconv.Atoi(r.URL.Query().Get("curah"))
	luas, _ := strconv.ParseFloat(r.URL.Query().Get("luas"), 64)

	kebutuhan := utils.HitungKebutuhanBibit(luas)

	var bibit models.Bibit
	err := database.DB.QueryRow(`
		SELECT id,nama,kualitas,stok FROM bibit
		WHERE LOWER(tanah)=LOWER($1) AND curah_hujan<=$2
		ORDER BY curah_hujan DESC LIMIT 1
	`, soil, rain).Scan(&bibit.ID, &bibit.Nama, &bibit.Kualitas, &bibit.Stok)

	if err != nil {
		utils.WriteJSON(w, 200, map[string]interface{}{
			"tanah":       soil,
			"curah_hujan": rain,
			"kebutuhan":   kebutuhan,
			"rekomendasi": "Tidak ada bibit yang cocok",
		})
		return
	}

	utils.WriteJSON(w, 200, map[string]interface{}{
		"tanah":       soil,
		"curah_hujan": rain,
		"kebutuhan":   kebutuhan,
		"rekomendasi": bibit,
	})
}

// LAPORAN HANDLER
func GetLaporan(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.UserKey).(map[string]interface{})
	username := ""
	if claims != nil {
		username = claims["username"].(string)
	}

	rows, err := database.DB.Query(`
		SELECT s.id, s.bibit_id, b.nama AS bibit_nama, s.user_id, u.username AS user_nama,
			   s.tipe, s.jumlah, s.created_at
		FROM stok_history s
		JOIN bibit b ON b.id=s.bibit_id
		JOIN users u ON u.id=s.user_id
		ORDER BY s.created_at DESC
	`)
	if err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	type rowData struct {
		SH        models.StokHistory
		BibitNama string
		UserNama  string
	}
	var rowsData []rowData
	for rows.Next() {
		var sh models.StokHistory
		var bibitNama, userNama string
		if err := rows.Scan(&sh.ID, &sh.BibitID, &bibitNama, &sh.UserID, &userNama, &sh.Tipe, &sh.Jumlah, &sh.CreatedAt); err != nil {
			utils.WriteError(w, 500, err.Error())
			return
		}
		rowsData = append(rowsData, rowData{sh, bibitNama, userNama})
	}

	result := utils.Map(rowsData, func(r rowData) map[string]interface{} {
		return utils.ToLaporanJSON(r.SH, r.BibitNama, r.UserNama)
	})

	totalMasuk := utils.Reduce(result, 0, func(acc int, r map[string]interface{}) int {
		if r["tipe"] == "masuk" {
			return acc + r["jumlah"].(int)
		}
		return acc
	})

	utils.WriteJSON(w, 200, map[string]interface{}{
		"laporan":     result,
		"total_masuk": totalMasuk,
		"user":        username,
	})
}

// USER HANDLER
func Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.WriteError(w, 400, "invalid body")
		return
	}

	validations := []utils.Validator{
		utils.ValidateNotEmpty(u.Username, "Username"),
		utils.ValidateNotEmpty(u.Password, "Password"),
		utils.ValidatePasswordLength(u.Password, 4),
	}
	for _, v := range validations {
		if ok, msg := v(); !ok {
			utils.WriteError(w, 400, msg)
			return
		}
	}

	_, err := database.DB.Exec(`INSERT INTO users (username,password,role) VALUES($1,$2,'user')`, u.Username, u.Password)
	if err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}

	utils.WriteJSON(w, 201, map[string]string{"message": "user registered"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		utils.WriteError(w, 400, "invalid body")
		return
	}

	var user models.User
	err := database.DB.QueryRow(`SELECT id,username,password,role FROM users WHERE username=$1`, loginData.Username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		utils.WriteError(w, 404, "user not found")
		return
	}

	if loginData.Password != user.Password {
		utils.WriteError(w, 401, "wrong password")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.WriteError(w, 500, "failed to generate token")
		return
	}

	utils.WriteJSON(w, 200, map[string]interface{}{
		"message": "login success",
		"token":   token,
		"user":    utils.ToUserJSON(user),
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`SELECT id,username,role,created_at FROM users ORDER BY id ASC`)
	if err != nil {
		utils.WriteError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt); err != nil {
			utils.WriteError(w, 500, err.Error())
			return
		}
		users = append(users, u)
	}

	result := utils.Map(users, utils.ToUserJSON)
	utils.WriteJSON(w, 200, result)
}
