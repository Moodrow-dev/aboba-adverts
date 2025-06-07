package handlers

import (
	"advert-server/database"
	"advert-server/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	MaxUploadSize  = 10 << 20
	UploadBasePath = "/uploads/"
)

func init() {
	var err error
	db, err = database.InitDB()
	if err != nil {
		panic(fmt.Errorf("database init failed: %w", err))
	}

	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		panic(fmt.Errorf("failed to create upload dir: %w", err))
	}
}

func AdvertHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAdverts(w, r)
	case http.MethodPost:
		createAdvert(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func AdvertDetailHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAdvert(w, r)
	case http.MethodPut:
		updateAdvert(w, r)
	case http.MethodDelete:
		deleteAdvert(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func checkOwnership(advertID, userID int) error {
	var ownerID int
	err := db.QueryRow("SELECT user_id FROM adverts WHERE id = ?", advertID).Scan(&ownerID)
	if err != nil {
		return err
	}
	if ownerID != userID {
		return errors.New("forbidden: not owner")
	}
	return nil
}

func createAdvert(w http.ResponseWriter, r *http.Request) {
	var advertReq models.AdvertRequest
	if err := json.NewDecoder(r.Body).Decode(&advertReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if advertReq.Title == "" || advertReq.Description == "" || advertReq.Price == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	userID := 1 // Temp - replace with auth

	result, err := db.Exec(
		"INSERT INTO adverts (title, description, price, category, user_id) VALUES (?, ?, ?, ?, ?)",
		advertReq.Title, advertReq.Description, advertReq.Price, advertReq.Category, userID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	advert := models.Advert{
		ID:          int(id),
		Title:       advertReq.Title,
		Description: advertReq.Description,
		Price:       advertReq.Price,
		Category:    advertReq.Category,
		UserID:      userID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(advert)
}

func getAdverts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, description, price, category, created_at FROM adverts ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var adverts []models.Advert
	for rows.Next() {
		var advert models.Advert
		if err := rows.Scan(&advert.ID, &advert.Title, &advert.Description, &advert.Price, &advert.Category, &advert.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		adverts = append(adverts, advert)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adverts)
}

func getAdvert(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/adverts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var advert models.Advert
	err = db.QueryRow(
		"SELECT id, title, description, price, category, created_at, user_id FROM adverts WHERE id = ?", id,
	).Scan(&advert.ID, &advert.Title, &advert.Description, &advert.Price, &advert.Category, &advert.CreatedAt, &advert.UserID)

	if err == sql.ErrNoRows {
		http.Error(w, "Not found", http.StatusNotFound)
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

func updateAdvert(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/adverts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var advertReq models.AdvertRequest
	if err := json.NewDecoder(r.Body).Decode(&advertReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if advertReq.Title == "" || advertReq.Description == "" || advertReq.Price == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	userID := 1 // Temp - replace with auth
	if err := checkOwnership(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	_, err = db.Exec(
		"UPDATE adverts SET title = ?, description = ?, price = ?, category = ? WHERE id = ?",
		advertReq.Title, advertReq.Description, advertReq.Price, advertReq.Category, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteAdvert(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/adverts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userID := 1 // Temp - replace with auth
	if err := checkOwnership(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	_, err = db.Exec("DELETE FROM adverts WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
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

	ext := filepath.Ext(header.Filename)
	newFilename := fmt.Sprintf("adv_%s_%d%s", advertID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(UploadDir, newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Save failed", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Save failed", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(
		"INSERT INTO photos (advert_id, url) VALUES (?, ?)",
		advertID, UploadBasePath+newFilename,
	)
	if err != nil {
		os.Remove(filePath)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"url": UploadBasePath + newFilename,
	})
}
