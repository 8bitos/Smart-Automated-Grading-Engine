package models

import "github.com/golang-jwt/jwt/v5"

// Digunakan untuk request body saat login
type LoginCredentials struct {
	Identifier string `json:"identifier"` // Bisa email atau username
	Password   string `json:"password"`
}

// Representasi data pengguna
type User struct {
	ID          string `json:"id,omitempty"`
	NamaLengkap string `json:"nama_lengkap"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	Peran       string `json:"peran"`
}

// Representasi data kelas
type Class struct {
	ID          string `json:"id,omitempty"`
	GuruID      string `json:"guru_id"`
	NamaKelas   string `json:"nama_kelas"`
	Deskripsi   string `json:"deskripsi,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// Payload untuk JWT
type Claims struct {
	UserID string `json:"user_id"`
	Peran  string `json:"peran"`
	jwt.RegisteredClaims
}
