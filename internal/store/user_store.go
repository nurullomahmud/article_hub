package store

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
	GetUserByID(id int64) (*User, error)
	UpdateUser(*User) (*User, error)
	DeleteUser(id int64) error
}

func (pg *PostgresUserStore) CreateUser(user *User) (*User, error) {
	query := `
	INSERT INTO users(email, hashed_password)
	VALUES ($1, $2)
	RETURNING id
	`
	err := pg.db.QueryRow(query, user.Email, user.HashedPassword).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	// handle validations and edge cases later
	return user, nil
}

func (pg *PostgresUserStore) GetUserByID(id int64) (*User, error) {
	user := &User{}
	query := `
	SELECT id, email, hashed_password
	FROM users
	WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.HashedPassword)
	if err == sql.ErrNoRows {
		return nil, nil // so we can catch not found easly in api layer
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pg *PostgresUserStore) UpdateUser(user *User) (*User, error) {
	updatedUser := &User{}
	query := `
	UPDATE users
	SET email = $1, hashed_password = $2
	WHERE id = $3
	RETURNING id, email, hashed_password
	`
	err := pg.db.QueryRow(query, user.Email, user.HashedPassword).Scan(&updatedUser.ID, &updatedUser.Email, &updatedUser.HashedPassword)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (pg *PostgresUserStore) DeleteUser(id int64) error {
	query := `
	DELETE FROM users WHERE id = $1
	`
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d: %w", id, sql.ErrNoRows)
	}

	return nil
}
