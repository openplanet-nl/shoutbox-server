package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type ResponseMessage struct {
	AccountID   string `json:"account_id"`
	DisplayName string `json:"display_name"`
	Message     string `json:"message"`
	Time        string `json:"time"`
}

type ResponseList struct {
	Items []ResponseMessage `json:"items"`
}

func apiList(w http.ResponseWriter, r *http.Request) {
	res := ResponseList{
		Items: make([]ResponseMessage, 0),
	}

	var messages []Message
	gDatabase.Order("time DESC").Limit(25).Find(&messages)
	for _, msg := range messages {
		res.Items = append(res.Items, ResponseMessage{
			AccountID:   msg.AccountID,
			DisplayName: msg.DisplayName,
			Message:     msg.Text,
			Time:        msg.Time.Format(time.RFC1123Z),
		})
	}

	resBytes, _ := json.Marshal(&res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resBytes)
}
