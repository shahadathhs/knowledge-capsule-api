package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"knowledge-capsule-api/models"
	"knowledge-capsule-api/store"
	"knowledge-capsule-api/utils"
)

var TopicStore = &store.TopicStore{FileStore: store.FileStore[models.Topic]{FilePath: "data/topics.json"}}

func TopicHandler(w http.ResponseWriter, r *http.Request) {
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
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}
