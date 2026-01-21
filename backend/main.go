package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// --- Konfigurasi Global ---
var jwtKey = []byte("kunci_rahasia_super_aman_yang_harus_diganti") // Pindahkan ke .env
var db *sql.DB

type contextKey string
const userClaimsKey = contextKey("userClaims")


// --- Structs ---
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID          string `json:"id,omitempty"`
	NamaLengkap string `json:"nama_lengkap"`
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"` // Password tidak akan di-JSON-kan saat marshaling jika kosong
	Peran       string `json:"peran"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Peran  string `json:"peran"`
	jwt.RegisteredClaims
}


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

	// Middleware Global
	r.Use(corsMiddleware)

	// Rute Publik
	r.HandleFunc("/api/hello", helloHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/auth/register", registerHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST", "OPTIONS")

	// Rute Admin (dilindungi oleh middleware)
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(jwtMiddleware, adminRequired)
	adminRouter.HandleFunc("/teachers", createTeacherHandler).Methods("POST", "OPTIONS")
	adminRouter.HandleFunc("/teachers", getTeachersHandler).Methods("GET", "OPTIONS") // NEW HANDLER
	adminRouter.HandleFunc("/teachers/{id}", deleteTeacherHandler).Methods("DELETE", "OPTIONS")
	
	log.Println("Go backend server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Tidak dapat memulai server: %s\n", err)
	}
}


// --- Handlers Publik---

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"message": "Request body tidak valid"}`, http.StatusBadRequest)
		return
	}

	var storedUser User
	var storedPasswordHash string
	query := `SELECT id, nama_lengkap, email, password, peran FROM users WHERE email=$1`
	err := db.QueryRow(query, creds.Email).Scan(&storedUser.ID, &storedUser.NamaLengkap, &storedUser.Email, &storedPasswordHash, &storedUser.Peran)
	if err != nil {
		http.Error(w, `{"message": "Email atau password salah"}`, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(creds.Password)); err != nil {
		http.Error(w, `{"message": "Email atau password salah"}`, http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: storedUser.ID,
		Peran:  storedUser.Peran,
		RegisteredClaims: jwt.RegisteredClaims{ ExpiresAt: jwt.NewNumericDate(expirationTime) },
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, `{"message": "Gagal membuat token"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"message": "Request body tidak valid"}`, http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" || user.NamaLengkap == "" {
		http.Error(w, `{"message": "Data tidak lengkap"}`, http.StatusBadRequest)
		return
	}
    // Peran otomatis student dari frontend, jadi validasi peran dihapus
    user.Peran = "student"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	query := `INSERT INTO users (nama_lengkap, email, password, peran) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.QueryRow(query, user.NamaLengkap, user.Email, string(hashedPassword), user.Peran).Scan(&user.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			http.Error(w, `{"message": "Email sudah terdaftar"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"message": "Gagal menyimpan pengguna"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registrasi berhasil!", "userID": user.ID})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello from Go backend!"})
}

// --- Handlers Admin ---

func createTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"message": "Request body tidak valid"}`, http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" || user.NamaLengkap == "" {
		http.Error(w, `{"message": "Data guru tidak lengkap"}`, http.StatusBadRequest)
		return
	}
    user.Peran = "teacher"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	query := `INSERT INTO users (nama_lengkap, email, password, peran) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.QueryRow(query, user.NamaLengkap, user.Email, string(hashedPassword), user.Peran).Scan(&user.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			http.Error(w, `{"message": "Email sudah terdaftar"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"message": "Gagal menyimpan guru"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Akun guru berhasil dibuat!", "userID": user.ID})
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, nama_lengkap, email, peran FROM users WHERE peran = 'teacher'`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Gagal mendapatkan daftar guru: %v", err)
		http.Error(w, `{"message": "Gagal mendapatkan daftar guru"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teachers []User
	for rows.Next() {
		var teacher User
		if err := rows.Scan(&teacher.ID, &teacher.NamaLengkap, &teacher.Email, &teacher.Peran); err != nil {
			log.Printf("Gagal scan data guru: %v", err)
			http.Error(w, `{"message": "Gagal memproses data guru"}`, http.StatusInternalServerError)
			return
		}
		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error setelah iterasi baris: %v", err)
		http.Error(w, `{"message": "Gagal mengambil data guru"}`, http.StatusInternalServerError)
		return
	}

	if teachers == nil {
        teachers = []User{} // Mengembalikan array kosong daripada nil
    }

	json.NewEncoder(w).Encode(teachers)
}

func deleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, `{"message": "ID guru tidak ditemukan"}`, http.StatusBadRequest)
		return
	}

	query := "DELETE FROM users WHERE id = $1 AND peran = 'teacher'"
	result, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, `{"message": "Gagal menghapus guru"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"message": "Guru tidak ditemukan atau pengguna bukan guru"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Akun guru berhasil dihapus"})
}


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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		// gorilla/mux sudah menangani OPTIONS secara otomatis jika method di-list
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
