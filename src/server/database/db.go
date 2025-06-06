package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB инициализирует базу данных и создает таблицы, если они не существуют
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./adverts.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Создание таблицы объявлений
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS adverts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		price TEXT NOT NULL,
		category TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id INTEGER
	);
	CREATE TABLE photos (
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		advert_id INTEGER NOT NULL,
		url VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (advert_id) REFERENCES adverts(id) ON DELETE CASCADE
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create adverts table: %v", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}
