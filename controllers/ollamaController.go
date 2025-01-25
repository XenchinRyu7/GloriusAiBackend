package controllers

import (
	"encoding/json"
	"gloriusaiapi/services"
	"io"
	"net/http"
)

type SetModelRequest struct {
	Name string `json:"name"`
}

func SetModel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SetModelRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Model name is required", http.StatusBadRequest)
		return
	}

	err = services.ActivateModel(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Model activated successfully"))
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req services.SendMessageRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Stream {
		_, err := services.SendMessageToModel(body, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response, err := services.SendMessageToModel(body, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
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
