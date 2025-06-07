package models

import (
	"time"
)

type Advert struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id,omitempty"`
	Photos      []Photo   `json:"photos,omitempty"`
}

type AdvertRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Category    string `json:"category"`
}

type Photo struct {
	ID       int    `json:"id"`
	AdvertID int    `json:"advert_id"`
	URL      string `json:"url"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}
