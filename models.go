package main

import "time"

type User struct {
	ID int `gorm:"primaryKey"`

	AccountID   string
	DisplayName string
	Secret      string
}

type Message struct {
	ID int `gorm:"primaryKey"`

	AccountID   string
	DisplayName string

	Time time.Time
	Text string
}
