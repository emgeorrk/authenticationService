package postgres

import (
	"authenticationService/internal/config"
	"authenticationService/internal/models"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storage config.Storage) (*Storage, error) {
	const op = "storage.postgres.NewStorage"

	connectString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		storage.User, storage.Password, storage.Database, storage.Host, storage.Port)

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
    	id UUID PRIMARY KEY,
    	name VARCHAR(255) NOT NULL,
    	email VARCHAR(255) UNIQUE NOT NULL);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: создать индекс на token_hash
	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS refresh_tokens (
    	id UUID PRIMARY KEY,
    	user_id UUID NOT NULL REFERENCES users(id),
    	refresh_token_hash VARCHAR(255) NOT NULL,
    	ip_address VARCHAR(39) NOT NULL,
	    status TEXT NOT NULL CHECK (status IN ('used', 'unused')),
    	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    	expires_at TIMESTAMP WITH TIME ZONE NOT NULL);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// TODO: подумать над возвращаемыми значениями

func (s *Storage) CreateUser(user *models.User) (string, error) {
	const op = "storage.postgres.CreateUser"

	id := uuid.New()

	stmt, err := s.db.Prepare("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(id.String(), user.Name, user.Email); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id.String(), nil
}

func (s *Storage) GetUserByID(id string) (*models.User, error) {
	const op = "storage.postgres.GetUserByID"

	stmt, err := s.db.Prepare("SELECT id, email FROM users WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User
	if err = stmt.QueryRow(id).Scan(&user.ID, &user.Email); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) CreateRefreshToken(refreshToken *models.RefreshToken) (string, error) {
	const op = "storage.postgres.CreateRefreshToken"

	stmt, err := s.db.Prepare("INSERT INTO refresh_tokens (id, user_id, refresh_token_hash, ip_address, status, created_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6, &7)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id := uuid.New()

	if _, err = stmt.Exec(id.String(), refreshToken.UserID, refreshToken.RefreshTokenHash, refreshToken.IPAddress, refreshToken.Status, refreshToken.CreatedAt, refreshToken.ExpiresAt); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id.String(), nil
}

func (s *Storage) GetRefreshTokenByID(id string) (*models.RefreshToken, error) {
	const op = "storage.postgres.GetRefreshTokenByID"

	stmt, err := s.db.Prepare("SELECT id, user_id, refresh_token_hash, ip_address, status, created_at, expires_at FROM refresh_tokens WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var refreshToken models.RefreshToken
	if err = stmt.QueryRow(id).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.RefreshTokenHash, &refreshToken.IPAddress, &refreshToken.Status, &refreshToken.CreatedAt, &refreshToken.ExpiresAt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &refreshToken, nil
}

func (s *Storage) UpdateRefreshTokenStatus(id, newStatus string) error {
	const op = "storage.postgres.UpdateRefreshTokenStatus"

	stmt, err := s.db.Prepare("UPDATE refresh_tokens SET status = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(newStatus, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
