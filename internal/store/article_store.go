package store

import "database/sql"

type Article struct {
	ID         int         `json:"id"`
	Title      string      `json:"title"`
	Image      string      `json:"image"`
	AuthorID   int         `json:"author_id"`
	Paragraphs []Paragraph `json:"paragraphs"`
}

type Paragraph struct {
	ID        int    `json:"id"`
	ArticleID int    `json:"article_id"`
	Headline  string `json:"headline"`
	Body      string `json:"body"`
	Order     int    `json:"order"`
}

type PostgresArticleStore struct {
	db *sql.DB
}

func NewPostgresArticleStore(db *sql.DB) *PostgresArticleStore {
	return &PostgresArticleStore{}
}

type ArticleStore interface {
	CreateArticle(*Article) (*Article, error)
	GetArticleByID(id int64) (*Article, error)
}
