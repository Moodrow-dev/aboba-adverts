package models

import (
	"time"
)

// Advert представляет структуру объявления
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

// AdvertRequest представляет запрос на создание/обновление объявления
type AdvertRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Category    string `json:"category"`
}

// Photo представляет фотографию объявления
type Photo struct {
	ID       int    `json:"id"`
	AdvertID int    `json:"advert_id"`
	URL      string `json:"url"`
}
