package api

import (
	"database/sql"
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

func (ah *ArticleHandler) HandleUpdateArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := utils.ReadIDParam(r)
	if err != nil {
		ah.logger.Printf("error converting id: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id parameter"})
		return
	}

	existingArticle, err := ah.articleStore.GetArticleByID(articleID)
	if err != nil {
		ah.logger.Printf("internal server error: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingArticle == nil {
		http.NotFound(w, r)
		return
	}

	var updatedArticle struct {
		Title      *string           `json:"title"`
		Image      *string           `json:"image"`
		AuthorID   *int              `json:"author_id"`
		Paragraphs []store.Paragraph `json:"paragraphs"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedArticle)
	if err != nil {
		ah.logger.Printf("error decoding payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid payload format"})
		return
	}

	if updatedArticle.Title != nil {
		existingArticle.Title = *updatedArticle.Title
	}
	if updatedArticle.Image != nil {
		existingArticle.Image = *updatedArticle.Image
	}
	if updatedArticle.Paragraphs != nil {
		existingArticle.Paragraphs = updatedArticle.Paragraphs
	}

	// later we implement some validations here

	err = ah.articleStore.UpdateArticle(existingArticle)
	if err != nil {
		ah.logger.Printf("error updating article: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"article": existingArticle})
}

func (ah *ArticleHandler) HandleDeleteArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := utils.ReadIDParam(r)
	if err != nil {
		ah.logger.Printf("error converting id param: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id data type"})
		return
	}

	article, err := ah.articleStore.GetArticleByID(articleID)
	if err != nil {
		ah.logger.Printf("error getting article: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if err == sql.ErrNoRows {
		http.Error(w, "article not found", http.StatusNotFound)
	}

	if err != nil {
		ah.logger.Printf("error getting article: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = ah.articleStore.DeleteArticle(article)

	if err == sql.ErrNoRows {
		http.Error(w, "article not found", http.StatusNotFound)
	}

	if err != nil {
		ah.logger.Printf("error deleting article: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "article deleted"})
}
