package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var gDatabase *gorm.DB

type ResponseError struct {
	Error string `json:"error"`
}

func sendError(w http.ResponseWriter, err string) {
	log.Printf("API error: %s\n", err)

	resBytes, _ := json.Marshal(ResponseError{
		Error: err,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(resBytes)
}

func main() {
	var err error

	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Unable to open config file: %s\n", err.Error())
		return
	}

	gDatabase, err = gorm.Open(sqlite.Open(viper.GetString("database.path")), &gorm.Config{})
	if err != nil {
		log.Printf("Unable to open sqlite database: %s\n", err.Error())
		return
	}

	gDatabase.AutoMigrate(&User{}, &Message{})

	http.HandleFunc("/auth", apiAuth)
	http.HandleFunc("/list", apiList)
	http.HandleFunc("/send", apiSend)
	http.HandleFunc("/users", apiUsers)

	log.Printf("Listening on port 8000\n")
	http.ListenAndServe(":8000", nil)
}
