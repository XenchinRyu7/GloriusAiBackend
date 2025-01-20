package controllers

import (
	"encoding/json"
	"gloriusaiapi/services"
	"gloriusaiapi/config"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var creds map[string]string
	json.NewDecoder(r.Body).Decode(&creds)

	user, token, err := services.AuthenticateUser(creds["username"], creds["password"])
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
    var userData map[string]string
    json.NewDecoder(r.Body).Decode(&userData)

    user, err := services.RegisterUser(config.DB, userData["name"], userData["email"], userData["password"])
    if err != nil {
        http.Error(w, "Failed to register user", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

