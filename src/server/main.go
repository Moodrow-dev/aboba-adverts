package main

import (
	"advert-server/handlers"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load("settings.env")
	log.SetOutput(os.Stdout)
	log.Println(fmt.Printf("Starting server on :%v", os.Getenv("PORT")))

	mux := http.NewServeMux()

	mux.HandleFunc("/api/adverts", handlers.AdvertHandler)
	mux.HandleFunc("/api/adverts/", handlers.AdvertDetailHandler)
	mux.HandleFunc("/api/upload-photo", handlers.UploadPhotoHandler)

	fs := http.FileServer(http.Dir("./uploads"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	handler := loggingMiddleware(enableCORS(mux))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), handler))
}
