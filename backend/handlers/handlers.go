package handlers

import (
	"encoding/json"
	"net/http"
	"sistem-skripsi/backend/models"
	"sistem-skripsi/backend/store"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type Server struct {
	router *mux.Router
	store  store.Store
}

func NewServer(router *mux.Router, store store.Store) *Server {
	return &Server{
		router: router,
		store:  store,
	}
}

func (s *Server) RegisterRoutes() {
	// Rute Publik
	s.router.HandleFunc("/api/auth/register", s.handleRegister).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/api/auth/login", s.handleLogin).Methods("POST", "OPTIONS")

	// Rute Guru & Admin
	teacherRouter := s.router.PathPrefix("/api").Subrouter()
	teacherRouter.Use(JWTMiddleware, TeacherOrAdminRequired)
	teacherRouter.HandleFunc("/classes", s.handleCreateClass).Methods("POST", "OPTIONS")
	teacherRouter.HandleFunc("/classes", s.handleGetClasses).Methods("GET", "OPTIONS")

	// Rute Admin
	adminRouter := s.router.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(JWTMiddleware, AdminRequired)
	adminRouter.HandleFunc("/teachers", s.handleCreateTeacher).Methods("POST", "OPTIONS")
	adminRouter.HandleFunc("/teachers", s.handleGetTeachers).Methods("GET", "OPTIONS")
	adminRouter.HandleFunc("/teachers/{id}", s.handleDeleteTeacher).Methods("DELETE", "OPTIONS")
}


// --- Handlers ---

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var creds models.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Request body tidak valid"})
		return
	}

	user, storedPasswordHash, err := s.store.GetUserByIdentifier(creds.Identifier)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Identifier atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(creds.Password)); err != nil {
		WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Identifier atau password salah"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID: user.ID,
		Peran:  user.Peran,
		RegisteredClaims: jwt.RegisteredClaims{ ExpiresAt: jwt.NewNumericDate(expirationTime) },
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat token"})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	s.handleUserCreation(w, r, "student")
}

func (s *Server) handleCreateTeacher(w http.ResponseWriter, r *http.Request) {
	s.handleUserCreation(w, r, "teacher")
}

func (s *Server) handleUserCreation(w http.ResponseWriter, r *http.Request, peran string) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Request body tidak valid"})
		return
	}

	validate := validator.New()
    if err := validate.Struct(user); err != nil {
        WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Data tidak lengkap atau tidak valid"})
        return
    }
	
	user.Peran = peran
	if err := s.store.CreateUser(&user); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "users_username_key":
				WriteJSON(w, http.StatusConflict, map[string]string{"message": "Username telah digunakan"})
				return
			case "users_email_key":
				WriteJSON(w, http.StatusConflict, map[string]string{"message": "Email sudah terdaftar"})
				return
			}
		}
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menyimpan pengguna"})
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{"message": "Akun berhasil dibuat!", "userID": user.ID})
}

func (s *Server) handleGetTeachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := s.store.GetTeachers()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mendapatkan daftar guru"})
		return
	}
	WriteJSON(w, http.StatusOK, teachers)
}

func (s *Server) handleDeleteTeacher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rowsAffected, err := s.store.DeleteUserByID(id, "teacher")
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus guru"})
		return
	}
	if rowsAffected == 0 {
		WriteJSON(w, http.StatusNotFound, map[string]string{"message": "Guru tidak ditemukan"})
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"message": "Akun guru berhasil dihapus"})
}

func (s *Server) handleCreateClass(w http.ResponseWriter, r *http.Request) {
	var classData models.Class
	if err := json.NewDecoder(r.Body).Decode(&classData); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Request body tidak valid"})
		return
	}

	claims, _ := r.Context().Value(userClaimsKey).(*models.Claims)
	classData.GuruID = claims.UserID

	if err := s.store.CreateClass(&classData); err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat kelas"})
		return
	}
	
	WriteJSON(w, http.StatusCreated, classData)
}

func (s *Server) handleGetClasses(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(userClaimsKey).(*models.Claims)
	
	classes, err := s.store.GetClassesByTeacherID(claims.UserID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil daftar kelas"})
		return
	}

	WriteJSON(w, http.StatusOK, classes)
}
