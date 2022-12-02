package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

// Response type from the Openplanet backend
type OpenplanetResponseAuth struct {
	Error string `json:"error"`

	AccountID   string `json:"account_id"`
	DisplayName string `json:"display_name"`
	TokenTime   int64  `json:"token_time"`
}

// Authentication response that we send to the client
type ResponseAuth struct {
	AccountID   string `json:"account_id"`
	DisplayName string `json:"display_name"`
	Secret      string `json:"secret"`
}

func apiAuth(w http.ResponseWriter, r *http.Request) {
	token := r.PostFormValue("t")

	// Prepare data to send to the Openplanet backend
	params := url.Values{}
	params.Set("token", token)
	params.Set("secret", viper.GetString("auth.secret"))

	// Send data to Openplanet backend for token validation
	body := bytes.NewReader([]byte(params.Encode()))
	res, err := http.Post(viper.GetString("auth.base")+"/api/auth/validate", "application/x-www-form-urlencoded", body)
	if err != nil {
		sendError(w, "Couldn't communicate with Openplanet backend: "+err.Error())
		return
	}

	// Read the response from the server
	resBytes, _ := io.ReadAll(res.Body)
	resAuth := OpenplanetResponseAuth{}
	json.Unmarshal(resBytes, &resAuth)

	// If there was an error, we must reject the token
	if resAuth.Error != "" {
		sendError(w, "Openplanet backend error: "+resAuth.Error)
		return
	}

	// If we get here, the token can be considered valid as long as the account ID is expected
	if resAuth.AccountID == "" {
		sendError(w, "Unexpected account ID")
		return
	}

	// Find an existing user with the authenticated account ID
	user := User{}
	gDatabase.Where("account_id", resAuth.AccountID).Find(&user)

	// If the user is not yet in the database, create it here
	if user.ID == 0 {
		user.AccountID = resAuth.AccountID

		// Generate a secret
		b := make([]byte, 48)
		rand.Read(b)
		user.Secret = base64.URLEncoding.EncodeToString(b)
	}

	// Update the display name in case it has changed
	user.DisplayName = resAuth.DisplayName

	// Insert or update the user in the database
	gDatabase.Save(&user)

	// Send response to the client including our user secret
	authResponse, _ := json.Marshal(ResponseAuth{
		AccountID:   user.AccountID,
		DisplayName: user.DisplayName,
		Secret:      user.Secret,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(authResponse)
}
