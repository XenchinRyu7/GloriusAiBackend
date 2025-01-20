package controllers

import (
	"encoding/json"
	"gloriusaiapi/config"
	"gloriusaiapi/models"
	"gloriusaiapi/services"
	"net/http"
)

type MessageResponse struct {
	Message  *models.Message  `json:"message"`
	Contexts []models.Context `json:"contexts"`
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var messageData map[string]string
	json.NewDecoder(r.Body).Decode(&messageData)

	userID := uint(1)
	message, err := services.CreateMessage(config.DB, userID, messageData["content"], messageData["response"])
	if err != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}

	contexts, err := services.GetContexts(config.DB, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve contexts", http.StatusInternalServerError)
		return
	}

	response := MessageResponse{
		Message:  message,
		Contexts: contexts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	userID := uint(1)

	messages, err := services.GetMessages(config.DB, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	contexts, err := services.GetContexts(config.DB, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve contexts", http.StatusInternalServerError)
		return
	}

	response := struct {
		Messages []models.Message `json:"messages"`
		Contexts []models.Context `json:"contexts"`
	}{
		Messages: messages,
		Contexts: contexts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
