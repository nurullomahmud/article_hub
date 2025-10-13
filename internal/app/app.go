package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/NurulloMahmud/article_hub/internal/api"
	"github.com/NurulloMahmud/article_hub/internal/store"
)

type Application struct {
	Logger         *log.Logger
	ArticleHandler *api.ArticleHandler
	Db             *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	pgDb, err := store.Open()
	if err != nil {
		return nil, err
	}

	// handlers
	articleHandler := api.NewArticle()

	app := &Application{
		Logger:         logger,
		ArticleHandler: articleHandler,
		Db:             pgDb,
	}

	return app, nil
}
