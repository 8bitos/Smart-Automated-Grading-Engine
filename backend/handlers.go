package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// --- Handlers Publik---

func helloHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello from Go backend!"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"message": "Request body tidak valid"}`, http.StatusBadRequest)
		return
	}

	var storedUser User
	var storedPasswordHash string
	var username sql.NullString // Handle NULL username in DB

	query := `SELECT id, nama_lengkap, username, email, password, peran FROM users WHERE email=$1 OR username=$1`
	err := db.QueryRow(query, creds.Identifier).Scan(&storedUser.ID, &storedUser.NamaLengkap, &username, &storedUser.Email, &storedPasswordHash, &storedUser.Peran)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"message": "Identifier atau password salah"}`, http.StatusUnauthorized)
			return
		}
		log.Printf("Error query user: %v", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	
	if username.Valid {
		storedUser.Username = username.String
	} else {
        storedUser.Username = "" // Ensure it's an empty string if NULL
    }

	if err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(creds.Password)); err != nil {
		http.Error(w, `{"message": "Identifier atau password salah"}`, http.StatusUnauthorized)
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
	handleUserCreation(w, r, "student")
}

// --- Handlers Admin ---

func createTeacherHandler(w http.ResponseWriter, r *http.Request) {
	handleUserCreation(w, r, "teacher")
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, nama_lengkap, username, email, peran FROM users WHERE peran = 'teacher'`
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
		var username sql.NullString // Handle NULL username
		if err := rows.Scan(&teacher.ID, &teacher.NamaLengkap, &username, &teacher.Email, &teacher.Peran); err != nil {
			log.Printf("Gagal scan data guru: %v", err)
			http.Error(w, `{"message": "Gagal memproses data guru"}`, http.StatusInternalServerError)
			return
		}
		if username.Valid {
			teacher.Username = username.String
		} else {
            teacher.Username = "" // Ensure it's an empty string if NULL
        }
		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error setelah iterasi baris: %v", err)
		http.Error(w, `{"message": "Gagal mengambil data guru"}`, http.StatusInternalServerError)
		return
	}

	if teachers == nil {
        teachers = []User{}
    }
	json.NewEncoder(w).Encode(teachers)
}

func deleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	query := "DELETE FROM users WHERE id = $1 AND peran = 'teacher'"
	result, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, `{"message": "Gagal menghapus guru"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"message": "Guru tidak ditemukan"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Akun guru berhasil dihapus"})
}

// --- Logika Helper ---



func handleUserCreation(w http.ResponseWriter, r *http.Request, peran string) {

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {

		http.Error(w, `{"message": "Request body tidak valid"}`, http.StatusBadRequest)

		return

	}



	if user.Email == "" || user.Password == "" || user.NamaLengkap == "" || user.Username == "" {

		http.Error(w, `{"message": "Data tidak lengkap"}`, http.StatusBadRequest)

		return

	}

    user.Peran = peran



	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	query := `INSERT INTO users (nama_lengkap, username, email, password, peran) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := db.QueryRow(query, user.NamaLengkap, user.Username, user.Email, string(hashedPassword), user.Peran).Scan(&user.ID)

	

	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {

			switch pqErr.Constraint {

			case "users_username_key":

				http.Error(w, `{"message": "Username telah digunakan"}`, http.StatusConflict)

				return

			case "users_email_key":

				http.Error(w, `{"message": "Email sudah terdaftar"}`, http.StatusConflict)

				return

			}

		}

		http.Error(w, `{"message": "Gagal menyimpan pengguna"}`, http.StatusInternalServerError)

		return

	}



	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{"message": "Akun berhasil dibuat!", "userID": user.ID})

}





// --- Handlers Kelas ---



func createClassHandler(w http.ResponseWriter, r *http.Request) {

	var classData Class

	if err := json.NewDecoder(r.Body).Decode(&classData); err != nil {

		http.Error(w, `{"message": "Request body tidak valid"}`, http.StatusBadRequest)

		return

	}



	claims, ok := r.Context().Value(userClaimsKey).(*Claims)

	if !ok {

		http.Error(w, `{"message": "Klaim pengguna tidak ditemukan di konteks"}`, http.StatusInternalServerError)

		return

	}



	classData.GuruID = claims.UserID



	query := `INSERT INTO classes (guru_id, nama_kelas, deskripsi) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := db.QueryRow(query, classData.GuruID, classData.NamaKelas, classData.Deskripsi).Scan(&classData.ID, &classData.CreatedAt)

	if err != nil {

		log.Printf("Gagal membuat kelas: %v", err)

		http.Error(w, `{"message": "Gagal membuat kelas di database"}`, http.StatusInternalServerError)

		return

	}



	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(classData)

}



func getClassesHandler(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value(userClaimsKey).(*Claims)

	if !ok {

		http.Error(w, `{"message": "Klaim pengguna tidak ditemukan di konteks"}`, http.StatusInternalServerError)

		return

	}



	query := `SELECT id, guru_id, nama_kelas, deskripsi, created_at FROM classes WHERE guru_id = $1 ORDER BY created_at DESC`

	rows, err := db.Query(query, claims.UserID)

	if err != nil {

		log.Printf("Gagal mendapatkan kelas: %v", err)

		http.Error(w, `{"message": "Gagal mendapatkan daftar kelas"}`, http.StatusInternalServerError)

		return

	}

	defer rows.Close()



	var classes []Class

	for rows.Next() {

		var c Class

		var description sql.NullString

		if err := rows.Scan(&c.ID, &c.GuruID, &c.NamaKelas, &description, &c.CreatedAt); err != nil {

			log.Printf("Gagal scan data kelas: %v", err)

			http.Error(w, `{"message": "Gagal memproses data kelas"}`, http.StatusInternalServerError)

			return

		}

		if description.Valid {

			c.Deskripsi = description.String

		}

		classes = append(classes, c)

	}



	if classes == nil {

		classes = []Class{}

	}



	json.NewEncoder(w).Encode(classes)

}
