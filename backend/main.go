package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// --- Konfigurasi Global ---
var jwtKey = []byte("kunci_rahasia_super_aman_yang_harus_diganti") // Pindahkan ke .env
var db *sql.DB

type contextKey string
const userClaimsKey = contextKey("userClaims")


func main() {
	// Koneksi Database
	connStr := "postgres://user:password@localhost:5432/essay_scoring?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil { log.Fatal("Gagal membuka koneksi database:", err) }
	defer db.Close()
	if err = db.Ping(); err != nil { log.Fatal("Gagal terhubung ke database:", err) }
	log.Println("Berhasil terhubung ke database PostgreSQL!")

	// Router
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	// Rute Publik
	r.HandleFunc("/api/hello", helloHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/auth/register", registerHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST", "OPTIONS")

	// Rute Guru & Admin
	teacherRouter := r.PathPrefix("/api").Subrouter()
	teacherRouter.Use(jwtMiddleware, teacherOrAdminRequired)
	teacherRouter.HandleFunc("/classes", createClassHandler).Methods("POST", "OPTIONS")
	teacherRouter.HandleFunc("/classes", getClassesHandler).Methods("GET", "OPTIONS")

	// Rute Admin
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(jwtMiddleware, adminRequired)
	adminRouter.HandleFunc("/teachers", createTeacherHandler).Methods("POST", "OPTIONS")
	adminRouter.HandleFunc("/teachers", getTeachersHandler).Methods("GET", "OPTIONS")
	adminRouter.HandleFunc("/teachers/{id}", deleteTeacherHandler).Methods("DELETE", "OPTIONS")
	
	log.Println("Go backend server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Tidak dapat memulai server: %s\n", err)
	}
}
