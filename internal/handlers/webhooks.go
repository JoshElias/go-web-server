package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/JoshElias/go-web-server/internal"
	"github.com/JoshElias/go-web-server/internal/services"
)

var UserUpgradedEvent = "user.upgraded"
var POLKA_API_KEY = "f271c81ff7084ee5b99a5091b42d486e"

func WebhookPolka(w http.ResponseWriter, r *http.Request) {
	// Check for auth key
	apiKey := getApiKeyFromAuthHeader(r)
	if apiKey == "" || apiKey != POLKA_API_KEY {
		internal.RespondWithStatus(w, 401)
		return
	}

	// Process event
	decoder := json.NewDecoder(r.Body)
	event := internal.PolkaWebhookEvent{}
	err := decoder.Decode(&event)
	if err != nil {
		internal.RespondWithStatus(w, 500)
		return
	}

	if event.Event != UserUpgradedEvent {
		internal.RespondWithStatus(w, 204)
		return
	}

	userId := event.Data.UserID
	user, err := services.GetUserById(userId)
	if err != nil {
		if errors.Is(err, internal.UserNotFound) {
			internal.RespondWithStatus(w, 404)
			return

		}
		internal.RespondWithStatus(w, 500)
		return
	}
	patch := internal.UserEntity{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: true,
	}
	if _, err := services.UpdateUserById(userId, patch); err != nil {
		internal.RespondWithStatus(w, 500)
		return
	}
	internal.RespondWithStatus(w, 204)
	return
}

func getApiKeyFromAuthHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	return strings.TrimPrefix(authHeader, "ApiKey ")
}
