package handlers

import (
	"errors"
	"net/http"
	"time"

	"knowledge-capsule-api/models"
	"knowledge-capsule-api/store"
	"knowledge-capsule-api/utils"
)

var UserStore = &store.UserStore{FileStore: store.FileStore[models.User]{FilePath: "data/users.json"}}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodPost) {
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if !utils.ParseAndValidateBody(w, r, &req) {
		return
	}

	user, err := UserStore.AddUser(req.Name, req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, true, "User registered", map[string]string{
		"user_id": user.ID,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodPost) {
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !utils.ParseAndValidateBody(w, r, &req) {
		return
	}

	user, err := UserStore.FindByEmail(req.Email)
	if err != nil || user == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, time.Hour*24)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, "Login successful", map[string]string{
		"token": token,
	})
}
