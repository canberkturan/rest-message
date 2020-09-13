package main

// User model
type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Message model
type Message struct {
	ID       string `json:"id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
	SendDate string `json:"senddate"`
	ReadDate string `json:"readdate"`
}
