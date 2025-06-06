package handlers

import (
	"advert-server/database"
	"advert-server/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	_ "mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB

const (
	UploadDir      = "./uploads"
	MaxUploadSize  = 10 << 20 // 10 MB
	UploadBasePath = "/uploads/"
)

func init() {
	var err error
	db, err = database.InitDB()
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		panic(err)
	}
}

// ... (существующие обработчики AdvertHandler и AdvertDetailHandler)

// UploadPhotoHandler загружает фото для объявления
func UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, "File too large (max 10MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	advertID := r.FormValue("advert_id")
	if _, err := strconv.Atoi(advertID); err != nil {
		http.Error(w, "Invalid advert ID", http.StatusBadRequest)
		return
	}

	// Проверяем существование объявления
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM adverts WHERE id = ?)", advertID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Advert not found", http.StatusNotFound)
		return
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(header.Filename)
	newFilename := fmt.Sprintf("adv_%s_%d%s", advertID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(UploadDir, newFilename)

	// Сохраняем файл
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Сохраняем в БД
	_, err = db.Exec(
		"INSERT INTO photos (advert_id, url) VALUES (?, ?)",
		advertID, UploadBasePath+newFilename,
	)
	if err != nil {
		os.Remove(filePath) // Удаляем файл если не удалось сохранить в БД
		http.Error(w, "Failed to save photo info", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"url": UploadBasePath + newFilename,
	})
}

// getAdvert (обновлённая версия с фото)
func getAdvert(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/adverts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid advert ID", http.StatusBadRequest)
		return
	}

	var advert models.Advert
	err = db.QueryRow(
		"SELECT id, title, description, price, category, created_at, user_id FROM adverts WHERE id = ?",
		id,
	).Scan(&advert.ID, &advert.Title, &advert.Description, &advert.Price, &advert.Category, &advert.CreatedAt, &advert.UserID)

	if err == sql.ErrNoRows {
		http.Error(w, "Advert not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, url FROM photos WHERE advert_id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var photos []models.Photo
	for rows.Next() {
		var photo models.Photo
		if err := rows.Scan(&photo.ID, &photo.URL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		photo.AdvertID = id
		photos = append(photos, photo)
	}

	advert.Photos = photos
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(advert)
}
