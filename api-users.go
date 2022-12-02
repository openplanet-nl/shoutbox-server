package main

import (
	"encoding/json"
	"net/http"
)

type ResponseUser struct {
	ID          int    `json:"id"`
	AccountID   string `json:"account_id"`
	DisplayName string `json:"display_name"`
}

type ResponseUsers struct {
	Users []ResponseUser
}

func apiUsers(w http.ResponseWriter, r *http.Request) {
	res := ResponseUsers{
		Users: make([]ResponseUser, 0),
	}

	var users []User
	gDatabase.Find(&users)
	for _, user := range users {
		res.Users = append(res.Users, ResponseUser{
			ID:          user.ID,
			AccountID:   user.AccountID,
			DisplayName: user.DisplayName,
		})
	}

	resBytes, _ := json.Marshal(&res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resBytes)
}
