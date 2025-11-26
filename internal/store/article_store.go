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
	return &PostgresArticleStore{
		db: db,
	}
}

type ArticleStore interface {
	CreateArticle(*Article) (*Article, error)
	GetArticleByID(id int64) (*Article, error)
	UpdateArticle(article *Article) error
	DeleteArticle(article *Article) error
}

func (pg *PostgresArticleStore) CreateArticle(article *Article) (*Article, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO articles (title, image, author_id)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err = tx.QueryRow(query, article.Title, article.Image, article.AuthorID).Scan(&article.ID)
	if err != nil {
		return nil, err
	}

	for _, paragraph := range article.Paragraphs {
		paragraphQuery := `
		INSERT INTO (article_id, headline, body, order)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`
		err = tx.QueryRow(paragraphQuery, article.ID, paragraph.Headline, paragraph.Body, paragraph.Order).Scan(&paragraph.ID)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (pg *PostgresArticleStore) GetArticleByID(id int64) (*Article, error) {
	article := &Article{}
	query := `
	SELECT id, title, image, author_id
	FROM articles
	WHERE id = $1
	`

	err := pg.db.QueryRow(query, id).Scan(&article.ID, article.Title, article.Image, article.AuthorID)
	if err != nil {
		return nil, err
	}

	paragraphQuery := `
	SELECT id, article_id, headline, body, order
	FROM paragraphs
	WHERE article_id = $1
	ORDER BY order
	`
	rows, err := pg.db.Query(paragraphQuery, article.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var paragraph Paragraph
		err = rows.Scan(
			&paragraph.ID,
			&paragraph.ArticleID,
			&paragraph.Headline,
			&paragraph.Body,
			&paragraph.Order,
		)
		if err != nil {
			return nil, err
		}
		article.Paragraphs = append(article.Paragraphs, paragraph)
	}
	return article, nil
}

func (pg *PostgresArticleStore) UpdateArticle(article *Article) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateQuery := `
	update articles
	set title = $1, image = $2, author_id = $3
	where id = $4
	`
	result, err := tx.Exec(updateQuery, article.Title, article.Image, article.AuthorID)
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

	_, err = tx.Exec(`delete from paragraphs where article_id = $1`, article.ID)
	if err != nil {
		return err
	}

	for _, paragraph := range article.Paragraphs {
		query := `
		INSERT INTO paragraphs(article_id, headline, body, order)
		VALUES ($1, $2, $3, $4)
		`
		_, err = tx.Exec(query,
			article.ID,
			paragraph.Headline,
			paragraph.Body,
			paragraph.Order,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pg *PostgresArticleStore) DeleteArticle(article *Article) error {
	deleteQuery := `
	delete from articles where id = $1
	`
	result, err := pg.db.Exec(deleteQuery, article.ID)
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

func (pg *PostgresArticleStore) GetArticleOwner(articleID int64) (int, error) {
	var userID int
	query := `
	SELECT author_id
	FROM articles
	WHERE id = $1
	`

	err := pg.db.QueryRow(query, articleID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
