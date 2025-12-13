package utils

import (
	"errors"
	"fmt"
	"go-seed-api/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GENERIC MAP / FILTER / REDUCE
func Map[T any, R any](list []T, fn func(T) R) []R {
	result := make([]R, 0, len(list))
	for _, v := range list {
		result = append(result, fn(v))
	}
	return result
}

func Filter[T any](list []T, fn func(T) bool) []T {
	result := []T{}
	for _, v := range list {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func Reduce[T any, R any](list []T, init R, fn func(R, T) R) R {
	acc := init
	for _, v := range list {
		acc = fn(acc, v)
	}
	return acc
}

// PURE BUSINESS LOGIC
func HitungKebutuhanBibit(luas float64) int {
	if luas <= 0 {
		return 0
	}
	return 5 + HitungKebutuhanBibit(luas-0.1)
}

// Transformasi Bibit ke summary JSON
func ToSummary(b models.Bibit) map[string]interface{} {
	return map[string]interface{}{
		"id":       b.ID,
		"nama":     b.Nama,
		"kualitas": b.Kualitas,
		"stok":     b.Stok,
	}
}

// Transformasi StokHistory ke JSON
func ToLaporanJSON(sh models.StokHistory, bibitNama, userNama string) map[string]interface{} {
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

// Transformasi User ke JSON
func ToUserJSON(u models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"username":   u.Username,
		"role":       u.Role,
		"created_at": u.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// VALIDATOR HOF
type Validator func() (bool, string)

func ValidateNotEmpty(field, name string) Validator {
	return func() (bool, string) {
		if field == "" {
			return false, name + " tidak boleh kosong"
		}
		return true, ""
	}
}

func ValidatePositive(value int, name string) Validator {
	return func() (bool, string) {
		if value < 0 {
			return false, name + " harus positif"
		}
		return true, ""
	}
}

func ValidatePasswordLength(password string, min int) Validator {
	return func() (bool, string) {
		if len(password) < min {
			return false, fmt.Sprintf("Password minimal %d karakter", min)
		}
		return true, ""
	}
}

// HELPER
func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

// JWT GENERATOR
var jwtKey = []byte("secretkey123")

func GenerateToken(userID int, username, role string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token invalid")
}
