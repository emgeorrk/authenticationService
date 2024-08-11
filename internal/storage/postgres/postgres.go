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
    	email VARCHAR(255) UNIQUE NOT NULL,
    	max_active_token_pairs INT NOT NULL DEFAULT 5,
    	access_token_lifetime_minutes INT NOT NULL DEFAULT 60,
    	refresh_token_lifetime_minutes INT NOT NULL DEFAULT 129600,
	    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT check_token_lifetime CHECK (access_token_lifetime_minutes <= users.refresh_token_lifetime_minutes));
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS tokens (
    	jti UUID PRIMARY KEY,
    	user_id UUID NOT NULL REFERENCES users(id),
    	refresh_token_hash VARCHAR(255) NOT NULL,
    	ip_address VARCHAR(39) NOT NULL,
	    refresh_token_status TEXT NOT NULL CHECK (refresh_token_status IN ('used', 'unused')),
    	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    	access_token_expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    	refresh_token_expires_at TIMESTAMP WITH TIME ZONE NOT NULL);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`CREATE INDEX IF NOT EXISTS idx_refresh_token_hash ON tokens(refresh_token_hash);`)
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

func (s *Storage) CreateToken(token *models.Token) (string, error) {
	const op = "storage.postgres.CreateRefreshToken"

	stmt, err := s.db.Prepare(`
		INSERT INTO tokens (
                    jti,
                    user_id,
                    refresh_token_hash,
                    ip_address,
                    refresh_token_status,
                    access_token_expires_at,
                    refresh_token_expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, &7);
	`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id := uuid.New()

	if _, err = stmt.Exec(
		id.String(),
		token.UserID,
		token.RefreshTokenHash,
		token.IPAddress,
		token.RefreshTokenStatus,
		token.CreatedAt,
		token.AccessTokenExpiresAt,
		token.RefreshTokenExpiresAt,
	); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id.String(), nil
}

func (s *Storage) GetTokenByJTI(JTI string) (*models.Token, error) {
	const op = "storage.postgres.GetRefreshTokenByID"

	stmt, err := s.db.Prepare(`
		SELECT JTI,
		       user_id,
		       refresh_token_hash,
		       ip_address,
		       refresh_token_status,
		       created_at,
		       access_token_expires_at,
		       refresh_token_expires_at
		FROM tokens
		WHERE jti = $1;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var refreshToken models.Token
	if err = stmt.QueryRow(JTI).Scan(&refreshToken); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &refreshToken, nil
}

func (s *Storage) UpdateRefreshTokenStatus(JTI, newStatus string) error {
	const op = "storage.postgres.UpdateRefreshTokenStatus"

	stmt, err := s.db.Prepare("UPDATE tokens SET refresh_token_status = $1 WHERE jti = $2;")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(newStatus, JTI); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetTokensByUserId(userID string) ([]models.Token, error) {
	const op = "storage.postgres.GetTokensByUserId"

	stmt, err := s.db.Prepare(`
		SELECT JTI,
		       user_id,
		       refresh_token_hash,
		       ip_address,
		       refresh_token_status,
		       created_at,
		       access_token_expires_at,
		       refresh_token_expires_at
		FROM tokens
		WHERE user_id = $1;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var tokens []models.Token
	for rows.Next() {
		var token models.Token
		if err = rows.Scan(&token); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}
