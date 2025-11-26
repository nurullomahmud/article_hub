package store

import (
	"database/sql"
	"time"

	"github.com/NurulloMahmud/article_hub/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
	ConfirmToken(token string, scope string) (bool, error)
	GetUserByToken(token string) (*User, error)
}

func (t *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)
	`
	_, err = t.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `
	DELETE FROM tokens
	WHERE user_id = $1 AND scope = $2
	`
	_, err := t.db.Exec(query, userID, scope)
	return err
}

func (t *PostgresTokenStore) ConfirmToken(token string, scope string) (bool, error) {
	query := `
	SELECT 1 FROM tokens WHERE token = $1 AND scope = $2 AND expiry < NOW()
	`
	var exists int
	err := t.db.QueryRow(query, token, scope).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (t *PostgresTokenStore) GetUserByToken(token string) (*User, error) {
	var userID int
	user := &User{}
	query := `
	SELECT user_id
	FROM tokens
	WHERE token = $1
	`
	err := t.db.QueryRow(query, token).Scan(&userID)
	if err != nil {
		return nil, err
	}
	getUserQuery := `
	SELECT id, email, hashed_password
	FROM users
	WHERE id = $1
	`
	err = t.db.QueryRow(getUserQuery, userID).Scan(&user.ID, &user.Email, &user.HashedPassword.hash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
