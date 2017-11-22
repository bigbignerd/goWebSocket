package main

type Message struct {
	FromUser string `json:"fromUser"`
	ToUser   string `json:"toUser"`
	Content  string `json:"content"`
}
