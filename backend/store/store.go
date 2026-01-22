package store

import (
	"database/sql"
	"sistem-skripsi/backend/models"

	"golang.org/x/crypto/bcrypt"
)

// Interface untuk semua operasi database
type Store interface {
	// User methods
	CreateUser(user *models.User) error
	GetUserByIdentifier(identifier string) (*models.User, string, error)
	GetTeachers() ([]*models.User, error)
	DeleteUserByID(id string, role string) (int64, error)
	// Class methods
	CreateClass(class *models.Class) error
	GetClassesByTeacherID(teacherID string) ([]*models.Class, error)
}

// Implementasi Store untuk PostgreSQL
type PostgresStore struct {
	db *sql.DB
}

// Konstruktor untuk PostgresStore
func NewPostgresStore() (*PostgresStore, error) {
	connStr := "postgres://user:password@localhost:5432/essay_scoring?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

// --- Implementasi method untuk User ---

func (s *PostgresStore) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (nama_lengkap, username, email, password, peran) 
              VALUES ($1, $2, $3, $4, $5) 
              RETURNING id`
	
	return s.db.QueryRow(
		query,
		user.NamaLengkap,
		user.Username,
		user.Email,
		string(hashedPassword),
		user.Peran,
	).Scan(&user.ID)
}


func (s *PostgresStore) GetUserByIdentifier(identifier string) (*models.User, string, error) {
	var user models.User
	var storedPasswordHash string
	var username sql.NullString 

	query := `SELECT id, nama_lengkap, username, email, password, peran FROM users WHERE email=$1 OR username=$1`
	
	err := s.db.QueryRow(query, identifier).Scan(
		&user.ID,
		&user.NamaLengkap,
		&username,
		&user.Email,
		&storedPasswordHash,
		&user.Peran,
	)
	if err != nil {
		return nil, "", err
	}

	if username.Valid {
		user.Username = username.String
	}

	return &user, storedPasswordHash, nil
}

func (s *PostgresStore) GetTeachers() ([]*models.User, error) {
	query := `SELECT id, nama_lengkap, username, email, peran FROM users WHERE peran = 'teacher'`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []*models.User
	for rows.Next() {
		var teacher models.User
		var username sql.NullString
		if err := rows.Scan(&teacher.ID, &teacher.NamaLengkap, &username, &teacher.Email, &teacher.Peran); err != nil {
			return nil, err
		}
		if username.Valid {
			teacher.Username = username.String
		}
		teachers = append(teachers, &teacher)
	}
	return teachers, nil
}


func (s *PostgresStore) DeleteUserByID(id string, role string) (int64, error) {
	query := "DELETE FROM users WHERE id = $1 AND peran = $2"
	result, err := s.db.Exec(query, id, role)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// --- Implementasi method untuk Class ---

func (s *PostgresStore) CreateClass(class *models.Class) error {
	query := `INSERT INTO classes (guru_id, nama_kelas, deskripsi) 
              VALUES ($1, $2, $3) 
              RETURNING id, created_at`

	return s.db.QueryRow(
		query,
		class.GuruID,
		class.NamaKelas,
		class.Deskripsi,
	).Scan(&class.ID, &class.CreatedAt)
}

func (s *PostgresStore) GetClassesByTeacherID(teacherID string) ([]*models.Class, error) {
	query := `SELECT id, guru_id, nama_kelas, deskripsi, created_at FROM classes WHERE guru_id = $1 ORDER BY created_at DESC`
	rows, err := s.db.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []*models.Class
	for rows.Next() {
		var c models.Class
		if err := rows.Scan(&c.ID, &c.GuruID, &c.NamaKelas, &c.Deskripsi, &c.CreatedAt); err != nil {
			return nil, err
		}
		classes = append(classes, &c)
	}
	return classes, nil
}
