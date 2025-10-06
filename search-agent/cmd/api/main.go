package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/Nikhil1169/search-agent/internal/gemini" // <-- CHANGED
	"github.com/Nikhil1169/search-agent/internal/handler"
	"github.com/Nikhil1169/search-agent/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// --- Database Connection (no changes) ---
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	store, err := storage.NewStore(connStr)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	log.Println("Database connected successfully!")

	// --- Gemini AI Client ---
	aiClient := gemini.NewClient() // <-- CHANGED

	// --- Handler and Router Setup ---
	searchHandler := &handler.SearchHandler{
		AI:    aiClient,
		Store: store,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/search/chat", searchHandler.HandleSearch)

	port := "3004"
	log.Printf("Search agent server running on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
