package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// --- Middlewares ---

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"message": "Header otorisasi tidak ditemukan"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"message": "Token tidak valid"}`, http.StatusUnauthorized)
			return
		}
		
		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func adminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userClaimsKey).(*Claims)
		if !ok {
			http.Error(w, `{"message": "Tidak dapat memproses klaim pengguna"}`, http.StatusInternalServerError)
			return
		}

		if claims.Peran != "superadmin" {
			http.Error(w, `{"message": "Akses ditolak: Memerlukan hak akses superadmin"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func teacherOrAdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userClaimsKey).(*Claims)
		if !ok {
			http.Error(w, `{"message": "Tidak dapat memproses klaim pengguna"}`, http.StatusInternalServerError)
			return
		}

		if claims.Peran != "teacher" && claims.Peran != "superadmin" {
			http.Error(w, `{"message": "Akses ditolak: Memerlukan hak akses guru atau superadmin"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
