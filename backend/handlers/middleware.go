package handlers

import (
	"context"
	"net/http"
	"sistem-skripsi/backend/models"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const userClaimsKey = contextKey("userClaims")

// NOTE: Kunci ini harus dipindahkan ke environment variable di production!
var jwtKey = []byte("kunci_rahasia_super_aman_yang_harus_diganti")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Header otorisasi tidak ditemukan"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token tidak valid"})
			return
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userClaimsKey).(*models.Claims)
		if !ok {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Tidak dapat memproses klaim pengguna"})
			return
		}

		if claims.Peran != "superadmin" {
			WriteJSON(w, http.StatusForbidden, map[string]string{"message": "Akses ditolak: Memerlukan hak akses superadmin"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func TeacherOrAdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userClaimsKey).(*models.Claims)
		if !ok {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Tidak dapat memproses klaim pengguna"})
			return
		}

		if claims.Peran != "teacher" && claims.Peran != "superadmin" {
			WriteJSON(w, http.StatusForbidden, map[string]string{"message": "Akses ditolak: Memerlukan hak akses guru atau superadmin"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
