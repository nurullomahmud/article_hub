package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NurulloMahmud/article_hub/internal/store"
	"github.com/NurulloMahmud/article_hub/internal/utils"
)

type ArticleHandler struct {
	articleStore store.ArticleStore
	logger       *log.Logger
}

func NewArticle(articleStore store.ArticleStore, logger *log.Logger) *ArticleHandler {
	return &ArticleHandler{
		articleStore: articleStore,
		logger:       logger,
	}
}

func (ah *ArticleHandler) HandleGetArticleByID(w http.ResponseWriter, r *http.Request) {
	articleID, err := utils.ReadIDParam(r)
	if err != nil {
		ah.logger.Printf("Error reading id param: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid article id"})
		return
	}

	article, err := ah.articleStore.GetArticleByID(articleID)
	if err != nil {
		ah.logger.Printf("error getting article by id: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"article": article})
}

func (ah *ArticleHandler) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	var article store.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		ah.logger.Printf("error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	article.AuthorID = 1 // handle this later
	createdArticle, err := ah.articleStore.CreateArticle(&article)
	if err != nil {
		ah.logger.Printf("error creating article: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"article": createdArticle})
}
