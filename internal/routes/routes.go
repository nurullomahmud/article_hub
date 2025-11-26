package routes

import (
	"github.com/NurulloMahmud/article_hub/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// article endpoints
	r.Get("/articles/{id}", app.ArticleHandler.HandleGetArticleByID)
	r.Post("/articles/", app.ArticleHandler.HandleCreateArticle)
	r.Put("/articles/{id}", app.ArticleHandler.HandleUpdateArticle)
	r.Delete("/articles/{id}", app.ArticleHandler.HandleDeleteArticle)

	// user endpoints
	r.Post("/register/", app.UserHandler.HandleRegister)

	// password change and reset endpoints
	r.Post("/users/{id}/password-change/", app.UserHandler.HandlePasswordChange)
	r.Post("/tokens/{email}/password-reset-request/", app.TokenHandler.HandlePasswordResetRequestToken)
	r.Post("/tokens/{token}/password-reset/", app.TokenHandler.HandlePasswordReset)

	return r
}
