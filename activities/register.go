package activities

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/liuminhaw/activitist/rand"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type RegisterService struct {
	DB            *sql.DB
	BytesPerToken int
	Duration      time.Duration
}

type Register struct {
	ID        int
	UserID    string
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

func (service RegisterService) TokenCreate(userID string) (*Register, error) {
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(service.BytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("token create: %w", err)
	}
	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	register := Register{
		UserID:    userID,
		Token:     token,
		TokenHash: service.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	_, err = service.DB.Exec(`
		INSERT INTO registrations (token_hash, expires, line_id)
		VALUES ($1, $2, $3) ON CONFLICT (line_id) 
		DO UPDATE 
		SET token_hash = $1, expires = $2
	`, register.TokenHash, register.ExpiresAt, register.UserID)
	if err != nil {
		return nil, fmt.Errorf("token create: %w", err)
	}

	return &register, nil
}

func (service RegisterService) Register(userId, token string) error {
	var register Register

	tokenHash := service.hash(token)
	row := service.DB.QueryRow(`
		SELECT id, line_id, expires
		FROM registrations
		WHERE token_hash = $1
	`, tokenHash)
	err := row.Scan(
		&register.ID, &register.UserID, &register.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}
	if time.Now().After(register.ExpiresAt) {
		return fmt.Errorf("register token expired: %v", token)
	}
	err = service.delete(register.ID)
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}
	_, err = service.DB.Exec(`
		INSERT INTO users (line_id)
		VALUES ($1)
	`, register.UserID)
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	return nil
}

func (service RegisterService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (service RegisterService) delete(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM registrations
		WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
