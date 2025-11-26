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

func (h *UserHandler) HandlePasswordChange(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("invalid user id format: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid param id"})
	}

	user, err := h.userStore.GetUserByID(userID)
	if err != nil {
		h.logger.Printf("Error getting user by id: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}

	var passChangeReq struct {
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		NewPasswordConfirm string `json:"new_password_confirm"`
	}

	err = json.NewDecoder(r.Body).Decode(&passChangeReq)
	if err != nil {
		h.logger.Printf("Error decoding password change request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err})
		return
	}
	if passChangeReq.NewPassword != passChangeReq.NewPasswordConfirm {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "passwords do not match"})
		return
	}
	if l := len(passChangeReq.NewPassword); l < 6 || l > 32 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "password length should be between 6 and 32"})
		return
	}
	passwordMatched, err := user.HashedPassword.Matches(passChangeReq.NewPassword)
	if err != nil {
		h.logger.Printf("Error checking password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if !passwordMatched {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "old password is not confirmed"})
		return
	}

	err = user.HashedPassword.Set(passChangeReq.NewPassword)
	if err != nil {
		h.logger.Printf("error changing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = h.userStore.UpdateUser(user)
	if err != nil {
		h.logger.Printf("Error updating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "password changed"})
}
