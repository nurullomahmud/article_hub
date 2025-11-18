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
}

func (t *PostgresTokenStore) CreateToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
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
