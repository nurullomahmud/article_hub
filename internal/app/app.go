package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/NurulloMahmud/article_hub/internal/api"
	"github.com/NurulloMahmud/article_hub/internal/middleware"
	"github.com/NurulloMahmud/article_hub/internal/store"
	"github.com/NurulloMahmud/article_hub/migrations"
)

type Application struct {
	Logger         *log.Logger
	ArticleHandler *api.ArticleHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     middleware.UserMiddleware
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
	tokenStore := store.NewPostgresTokenStore(pgDb)

	// handlers
	articleHandler := api.NewArticle(articleStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middlewareHandler := middleware.UserMiddleware{TokenStore: tokenStore}

	app := &Application{
		Logger:         logger,
		ArticleHandler: articleHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		Middleware:     middlewareHandler,
		Db:             pgDb,
	}

	return app, nil
}
