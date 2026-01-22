package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sistem-skripsi/backend/models"
	"sistem-skripsi/backend/store"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type APIServer struct {
	listenAddr string
	store      store.Store
}

func NewAPIServer(listenAddr string, store store.Store) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.Use(CorsMiddleware)

	// Attach subrouters and handlers
	s.registerRoutes(router)

	// Start the server
	go func() {
		log.Println("Go backend server starting on", s.listenAddr)
		if err := http.ListenAndServe(s.listenAddr, router); err != nil {
			log.Fatalf("Tidak dapat memulai server: %s\n", err)
		}
	}()
}

func (s *APIServer) registerRoutes(router *mux.Router) {
	// Rute Publik
	router.HandleFunc("/api/auth/register", s.handleRegister).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", s.handleLogin).Methods("POST", "OPTIONS")

	// Rute Guru & Admin

teacherRouter := router.PathPrefix("/api").Subrouter()
teacherRouter.Use(JWTMiddleware, TeacherOrAdminRequired)
teacherRouter.HandleFunc("/classes", s.handleCreateClass).Methods("POST", "OPTIONS")
teacherRouter.HandleFunc("/classes", s.handleGetClasses).Methods("GET", "OPTIONS")

	// Rute Admin
	adminRouter := router.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(JWTMiddleware, AdminRequired)
	adminRouter.HandleFunc("/teachers", s.handleCreateTeacher).Methods("POST", "OPTIONS")
	adminRouter.HandleFunc("/teachers", s.handleGetTeachers).Methods("GET", "OPTIONS")
	adminRouter.HandleFunc("/teachers/{id}", s.handleDeleteTeacher).Methods("DELETE", "OPTIONS")
}

// --- Handlers ---
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
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

	// ... (JWT creation logic)
}

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	s.handleUserCreation(w, r, "student")
}

// ... (other handlers)
