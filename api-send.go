package main

import (
	"net/http"
	"strings"
	"time"
)

func apiSend(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Secret ") {
		sendError(w, "Invalid authorization")
		return
	}

	secret := strings.TrimPrefix(authHeader, "Secret ")

	var user User
	gDatabase.Where("secret = ?", secret).Find(&user)

	if user.ID == 0 {
		sendError(w, "Invalid authorization")
		return
	}

	msg := strings.TrimSpace(r.FormValue("message"))
	if len(msg) == 0 {
		sendError(w, "Message is empty")
		return
	}

	gDatabase.Save(&Message{
		AccountID:   user.AccountID,
		DisplayName: user.DisplayName,

		Time: time.Now(),
		Text: msg,
	})
}
