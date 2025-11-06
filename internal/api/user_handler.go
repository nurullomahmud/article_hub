package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NurulloMahmud/article_hub/internal/store"
	"github.com/NurulloMahmud/article_hub/internal/utils"
)

type UserHandler struct {
	Logger    *log.Logger
	UserStore store.UserStore
}

func NewUser(store store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		Logger:    logger,
		UserStore: store,
	}
}

func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadIDParam(r)
	if err != nil {
		uh.Logger.Printf("invalid id param: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id param"})
		return
	}

	user, err := uh.UserStore.GetUserByID(userID)
	if err != nil {
		uh.Logger.Printf("error getting user by id in API layer: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}

	if user == nil {
		http.NotFound(w, r)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}

func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser store.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		uh.Logger.Printf("error decoding payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	createdUser, err := uh.UserStore.CreateUser(&newUser)
	if err != nil {
		uh.Logger.Printf("error creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": createdUser})
}
