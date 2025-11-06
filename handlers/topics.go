package handlers

import (
	"encoding/json"
	"net/http"

	// "knowledge-capsule-api/middleware"
	"knowledge-capsule-api/models"
	"knowledge-capsule-api/store"
	"knowledge-capsule-api/utils"
)

var TopicStore = &store.TopicStore{FileStore: store.FileStore[models.Topic]{FilePath: "data/topics.json"}}

func TopicHandler(w http.ResponseWriter, r *http.Request) {
	// userID := r.Context().Value(middleware.UserContextKey).(string)

	switch r.Method {
	case http.MethodGet:
		topics, _ := TopicStore.GetAllTopics()
		utils.JSONResponse(w, http.StatusOK, true, "Topics fetched", topics)

	case http.MethodPost:
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		topic, err := TopicStore.AddTopic(req.Name, req.Description)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		utils.JSONResponse(w, http.StatusCreated, true, "Topic created", topic)

	default:
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, nil)
	}
}
