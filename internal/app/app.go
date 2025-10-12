package app

import (
	"log"
	"os"

	"github.com/NurulloMahmud/article_hub/internal/api"
)

type Application struct {
	Logger         *log.Logger
	ArticleHandler *api.ArticleHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// handlers
	articleHandler := api.NewArticle()

	app := &Application{
		Logger:         logger,
		ArticleHandler: articleHandler,
	}

	return app, nil
}
