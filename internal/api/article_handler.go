package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ArticleHandler struct{}

func NewArticle() *ArticleHandler {
	return &ArticleHandler{}
}

func (ah *ArticleHandler) HandleGetArticleByID(w http.ResponseWriter, r *http.Request) {
	paramsID := chi.URLParam(r, "id")
	if paramsID == "" {
		http.NotFound(w, r)
		return
	}

	idInt, err := strconv.ParseInt(paramsID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "This is the article by id %d\n", idInt)
}

func (ah *ArticleHandler) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating an article")
}

func (ah *ArticleHandler) HandleListArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of articles")
}
