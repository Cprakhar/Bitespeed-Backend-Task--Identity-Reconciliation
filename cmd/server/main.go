package main

import (
	"bitespeed-identity-reconciliation/internal/database"
	"bitespeed-identity-reconciliation/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	log.Println("Starting Bitespeed Identity Reconciliation Service...")

	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	identifyHandler := handlers.NewIdentifyHandler()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/identity", identifyHandler.Identify)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	
	log.Printf("Server is running on port %s", port)
	log.Println("Available endpoints:")
    log.Println("  GET  /health   - Health check")
    log.Println("  POST /identify - Identity reconciliation")
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}