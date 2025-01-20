package controllers

import (
	"encoding/json"
	"gloriusaiapi/services"
	"net/http"
	"gloriusaiapi/config"
)

func SaveContext(w http.ResponseWriter, r *http.Request) {
    var contextData map[string]string
    json.NewDecoder(r.Body).Decode(&contextData)

    userID := uint(1)
    context, err := services.CreateOrUpdateContext(config.DB, userID, contextData["key"], contextData["value"])
    if err != nil {
        http.Error(w, "Failed to save context", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(context)
}

func GetContexts(w http.ResponseWriter, r *http.Request) {
    userID := uint(1)
    contexts, err := services.GetContexts(config.DB, userID)
    if err != nil {
        http.Error(w, "Failed to retrieve contexts", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(contexts)
}

