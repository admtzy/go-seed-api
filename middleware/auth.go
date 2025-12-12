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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized: token invalid atau expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
