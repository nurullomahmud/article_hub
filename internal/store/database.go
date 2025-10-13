package store

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres dbname=article_hub password=postgres port=5432 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("error connecting to database -> %w", err)
	}

	fmt.Println("Connected to database...")
	return db, nil
}
