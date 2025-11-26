package api

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/NurulloMahmud/article_hub/internal/store"
	"github.com/NurulloMahmud/article_hub/internal/tokens"
	"github.com/NurulloMahmud/article_hub/internal/utils"
	"github.com/go-chi/chi/v5"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		userStore:  userStore,
		tokenStore: tokenStore,
		logger:     logger,
	}
}

func (th *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		th.logger.Printf("Error decoding request payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	user, err := th.userStore.GetUserByEmail(req.Email)
	if err != nil {
		th.logger.Printf("error getting user by email: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}
	// handle password match !!!!
	passwordMatched, err := user.HashedPassword.Matches(req.Password)
	if err != nil {
		th.logger.Printf("error matching password with hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordMatched {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := th.tokenStore.CreateNewToken(user.ID, time.Hour*24, tokens.ScopeAuth)
	if err != nil {
		th.logger.Printf("error creating token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"token": token})
}

func (th *TokenHandler) HandlePasswordResetRequestToken(w http.ResponseWriter, r *http.Request) {
	emailParam := chi.URLParam(r, "email")
	if emailParam == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "email is required"})
		return
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(emailParam) {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid email format"})
		return
	}

	user, err := th.userStore.GetUserByEmail(emailParam)
	if err != nil {
		th.logger.Printf("error getting user by email: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}

	newToken, err := th.tokenStore.CreateNewToken(user.ID, time.Minute*30, tokens.ScopePasswordReset)
	if err != nil {
		th.logger.Printf("error generating password reset request token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	th.logger.Printf("There is new token: %s", newToken)

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "token generated successfully"})
}

func (th *TokenHandler) HandlePasswordReset(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "token is required"})
		return
	}

	// check if the toke is still valid
	valid, err := th.tokenStore.ConfirmToken(token, tokens.ScopePasswordReset)
	if err != nil {
		
	}
}