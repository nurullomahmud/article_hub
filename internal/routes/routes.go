package routes

import (
	"github.com/NurulloMahmud/article_hub/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/articles/{id}", app.ArticleHandler.HandleGetArticleByID)
	r.Post("/articles/", app.ArticleHandler.HandleCreateArticle)
	
	return r
}