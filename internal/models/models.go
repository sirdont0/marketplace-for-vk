package models

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Login        string    `json:"login" binding:"required,min=3,max=50"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Ad struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" binding:"required,min=5,max=100"`
	Text      string    `json:"text" binding:"required,min=10,max=1000"`
	ImageURL  string    `json:"image_url" binding:"required,url"`
	Price     float64   `json:"price" binding:"required,min=0"`
	AuthorID  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AdResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	ImageURL  string    `json:"image_url"`
	Price     float64   `json:"price"`
	Author    string    `json:"author"`
	IsMine    bool      `json:"is_mine"`
	CreatedAt time.Time `json:"created_at"`
}
