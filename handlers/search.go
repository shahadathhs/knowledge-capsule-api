package handlers

import (
	"net/http"

	"knowledge-capsule-api/middleware"
	"knowledge-capsule-api/utils"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, nil)
		return
	}

	results, _ := CapsuleStore.SearchCapsules(userID, query)
	utils.JSONResponse(w, http.StatusOK, true, "Search results", results)
}
