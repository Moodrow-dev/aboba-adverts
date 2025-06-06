package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// --- Configuration ---
type Config struct {
	Port   string
	DBPath string
}

func loadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4242"
	}

	return Config{
		Port:   port,
		DBPath: "./data/adverts.db",
	}
}

// --- Data Models ---
type Advert struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id,omitempty"`
	Seller      *Seller   `json:"seller,omitempty"`
}

type AdvertRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       string   `json:"price"`
	Category    string   `json:"category"`
	Images      []string `json:"images,omitempty"`
}

type Seller struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// --- Repository ---
type AdvertRepository struct {
	db *sql.DB
}

func NewAdvertRepository(db *sql.DB) *AdvertRepository {
	return &AdvertRepository{db: db}
}

func (r *AdvertRepository) Create(advert *Advert, images []string) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`INSERT INTO adverts (title, description, price, category, user_id, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		advert.Title, advert.Description, advert.Price, advert.Category, advert.UserID, advert.CreatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("insert advert: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert ID: %w", err)
	}

	for _, img := range images {
		_, err := tx.Exec(
			`INSERT INTO advert_images (advert_id, image_url) VALUES (?, ?)`,
			id, img,
		)
		if err != nil {
			return 0, fmt.Errorf("insert image: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return int(id), nil
}

func (r *AdvertRepository) GetByID(id int) (*Advert, []string, error) {
	var advert Advert
	var seller Seller
	err := r.db.QueryRow(
		`SELECT a.id, a.title, a.description, a.price, a.category, a.created_at, a.user_id, 
			u.id, u.name, u.nickname, u.email
		FROM adverts a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.id = ?`, id,
	).Scan(&advert.ID, &advert.Title, &advert.Description, &advert.Price,
		&advert.Category, &advert.CreatedAt, &advert.UserID,
		&seller.ID, &seller.Name, &seller.Nickname, &seller.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Advert not found for ID: %d", id)
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("query advert: %w", err)
	}

	advert.Seller = &seller

	rows, err := r.db.Query(`SELECT image_url FROM advert_images WHERE advert_id = ?`, id)
	if err != nil {
		return nil, nil, fmt.Errorf("query images: %w", err)
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var img string
		if err := rows.Scan(&img); err != nil {
			return nil, nil, fmt.Errorf("scan image: %w", err)
		}
		images = append(images, img)
	}

	return &advert, images, nil
}

func (r *AdvertRepository) GetHotAdverts(limit int) ([]Advert, error) {
	rows, err := r.db.Query(
		`SELECT a.id, a.title, a.description, a.price, a.category, a.created_at, a.user_id,
			u.id, u.name, u.nickname, u.email
		FROM adverts a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC
		LIMIT ?`, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query hot adverts: %w", err)
	}
	defer rows.Close()

	var adverts []Advert
	for rows.Next() {
		var advert Advert
		var seller Seller
		err := rows.Scan(&advert.ID, &advert.Title, &advert.Description, &advert.Price,
			&advert.Category, &advert.CreatedAt, &advert.UserID,
			&seller.ID, &seller.Name, &seller.Nickname, &seller.Email)
		if err != nil {
			return nil, fmt.Errorf("scan advert: %w", err)
		}
		advert.Seller = &seller
		adverts = append(adverts, advert)
	}

	return adverts, nil
}

func (r *AdvertRepository) Delete(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete advert (images are automatically deleted due to ON DELETE CASCADE)
	res, err := tx.Exec(`DELETE FROM adverts WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete advert: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("advert with ID %d not found", id)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// --- API Handlers ---
type APIHandler struct {
	repo *AdvertRepository
}

func NewAPIHandler(repo *AdvertRepository) *APIHandler {
	return &APIHandler{repo: repo}
}

func (h *APIHandler) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	var req AdvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := validateAdvertRequest(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	advert := Advert{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		UserID:      1, // Temporary value, replace with authenticated user ID
		CreatedAt:   time.Now(),
	}

	id, err := h.repo.Create(&advert, req.Images)
	if err != nil {
		log.Printf("Create advert error: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to create advert")
		return
	}

	advert.ID = id
	respondJSON(w, http.StatusCreated, advert)
}

func (h *APIHandler) GetAdvert(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("Invalid advert ID: %s", mux.Vars(r)["id"])
		respondError(w, http.StatusBadRequest, "Invalid advert ID")
		return
	}

	advert, images, err := h.repo.GetByID(id)
	if err != nil {
		log.Printf("Get advert error for ID %d: %v", id, err)
		respondError(w, http.StatusInternalServerError, "Failed to get advert")
		return
	}

	if advert == nil {
		respondError(w, http.StatusNotFound, "Advert not found")
		return
	}

	type Response struct {
		Advert
		Images []string `json:"images"`
	}
	respondJSON(w, http.StatusOK, Response{*advert, images})
}

func (h *APIHandler) GetHotAdverts(w http.ResponseWriter, r *http.Request) {
	limit := 10 // Default limit for hot adverts
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	adverts, err := h.repo.GetHotAdverts(limit)
	if err != nil {
		log.Printf("Get hot adverts error: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get hot adverts")
		return
	}

	respondJSON(w, http.StatusOK, adverts)
}

func (h *APIHandler) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("Invalid advert ID: %s", mux.Vars(r)["id"])
		respondError(w, http.StatusBadRequest, "Invalid advert ID")
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		log.Printf("Delete advert error for ID %d: %v", id, err)
		if err.Error() == fmt.Sprintf("advert with ID %d not found", id) {
			respondError(w, http.StatusNotFound, "Advert not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete advert")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Advert deleted successfully"})
}

// --- Helper Functions ---
func validateAdvertRequest(req AdvertRequest) error {
	if req.Title == "" || len(req.Title) > 100 {
		return errors.New("title must be between 1 and 100 characters")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	if req.Price == "" {
		return errors.New("price is required")
	}
	return nil
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// --- Database Initialization ---
func initDB(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Create tables
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			nickname TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS adverts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			price TEXT NOT NULL,
			category TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			user_id INTEGER,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS advert_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			advert_id INTEGER NOT NULL,
			image_url TEXT NOT NULL,
			FOREIGN KEY (advert_id) REFERENCES adverts(id) ON DELETE CASCADE
		);
	`); err != nil {
		return nil, fmt.Errorf("create tables: %w", err)
	}

	// Insert a default user for testing
	_, err = db.Exec(`
		INSERT OR IGNORE INTO users (id, name, nickname, email) 
		VALUES (1, 'Test User', 'testuser', 'test@example.com')
	`)
	if err != nil {
		return nil, fmt.Errorf("insert default user: %w", err)
	}

	return db, nil
}

// --- Main Function ---
func main() {
	cfg := loadConfig()

	db, err := initDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	repo := NewAdvertRepository(db)
	handler := NewAPIHandler(repo)

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	// API Routes
	api.HandleFunc("/adverts", handler.CreateAdvert).Methods("POST")
	api.HandleFunc("/adverts", handler.GetHotAdverts).Methods("GET")
	api.HandleFunc("/adverts/{id}", handler.GetAdvert).Methods("GET")
	api.HandleFunc("/adverts/{id}", handler.DeleteAdvert).Methods("DELETE")

	// CORS setup
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      cors(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
