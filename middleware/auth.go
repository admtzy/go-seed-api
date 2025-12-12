package middleware

import (
	"context"
	"go-seed-api/utils"
	"net/http"
	"strings"
)

type key int

const UserKey key = 0

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", 401)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized: token invalid atau expired", 401)
			return
		}

		// Masukkan claims ke context supaya handler bisa akses user_id & username
		ctx := context.WithValue(r.Context(), UserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
