package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/NurulloMahmud/article_hub/internal/api"
	"github.com/NurulloMahmud/article_hub/internal/store"
	"github.com/NurulloMahmud/article_hub/migrations"
)

type Application struct {
	Logger         *log.Logger
	ArticleHandler *api.ArticleHandler
	UserHandler    *api.UserHandler
	Db             *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	pgDb, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// stores for tables
	articleStore := store.NewPostgresArticleStore(pgDb)
	userStore := store.NewPostgresUserStore(pgDb)

	// handlers
	articleHandler := api.NewArticle(articleStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)

	app := &Application{
		Logger:         logger,
		ArticleHandler: articleHandler,
		UserHandler:    userHandler,
		Db:             pgDb,
	}

	return app, nil
}
