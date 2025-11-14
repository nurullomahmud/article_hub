package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/NurulloMahmud/article_hub/internal/store"
	"github.com/NurulloMahmud/article_hub/internal/utils"
)

type registerUserRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		logger:    logger,
		userStore: userStore,
	}
}

func (r *registerUserRequest) validateBasic() error {
	// check email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}

	// check password validation
	if len(r.Password) < 6 || len(r.Password) > 32 {
		return errors.New("password length must be between 6 and 32")
	}
	if r.Password != r.PasswordConfirm {
		return errors.New("passwords do not match")
	}

	return nil
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding request payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = req.validateBasic()
	if err != nil {
		h.logger.Printf("invalid request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	existingUser, err := h.userStore.GetUserByEmail(req.Email)
	if err != nil {
		h.logger.Printf("error getting user by email: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if existingUser != nil {
		h.logger.Printf("duplicate email entry")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "email already exists"})
		return
	}

	user := &store.User{
		Email: req.Email,
	}
	err = user.HashedPassword.Set(req.Password)
	if err != nil {
		h.logger.Printf("error hashing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("error creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}
