package handlers

import (
	"encoding/json"
	"net/http"

	"knowledge-capsule-api/middleware"
	"knowledge-capsule-api/models"
	"knowledge-capsule-api/store"
	"knowledge-capsule-api/utils"
)

var CapsuleStore = &store.CapsuleStore{FileStore: store.FileStore[models.Capsule]{FilePath: "data/capsules.json"}}

func CapsuleHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)

	switch r.Method {
	case http.MethodGet:
		capsules, _ := CapsuleStore.GetCapsulesByUser(userID)
		utils.JSONResponse(w, http.StatusOK, true, "Capsules fetched", capsules)

	case http.MethodPost:
		var req struct {
			Title     string   `json:"title"`
			Content   string   `json:"content"`
			Topic     string   `json:"topic"`
			Tags      []string `json:"tags"`
			IsPrivate bool     `json:"is_private"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		capsule, err := CapsuleStore.AddCapsule(userID, req.Title, req.Content, req.Topic, req.Tags, req.IsPrivate)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		utils.JSONResponse(w, http.StatusCreated, true, "Capsule created", capsule)

	default:
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, nil)
	}
}
