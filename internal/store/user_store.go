package store

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, err
	}
	
	return true, nil
}

type User struct {
	ID             int      `json:"id"`
	Email          string   `json:"email"`
	HashedPassword password `json:"-"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByEmail(email string) (*User, error)
	UpdateUser(*User) error
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (email, hashed_password)
	VALUES ($1, $2)
	RETURNING id 
	`
	err := s.db.QueryRow(query, user.Email, user.HashedPassword.hash).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) GetUserByEmail(email string) (*User, error) {
	query := `
	SELECT id, email, hashed_password FROM users WHERE email = $1
	`
	user := &User{
		HashedPassword: password{},
	}
	err := s.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.HashedPassword.hash)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	updateUserQuery := `
	update users
	set email = $1
	where id = $2
	`
	result, err := s.db.Exec(updateUserQuery, user.Email, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
