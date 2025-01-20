package controllers

import (
	"encoding/json"
	"net/http"
	"gloriusaiapi/services"
)

func SetModel(w http.ResponseWriter, r *http.Request) {
	modelName := r.URL.Query().Get("name")
	if modelName == "" {
		http.Error(w, "Model name is required", http.StatusBadRequest)
		return
	}

	err := services.ActivateModel(modelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Model activated successfully"))
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := services.SendMessageToModel(request.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"response": response,
	})
}

func GetAllModels(w http.ResponseWriter, r *http.Request) {
	models, err := services.GetAllModels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}
