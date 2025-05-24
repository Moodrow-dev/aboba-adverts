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
}

// AdvertRequest представляет запрос на создание/обновление объявления
type AdvertRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Category    string `json:"category"`
}
