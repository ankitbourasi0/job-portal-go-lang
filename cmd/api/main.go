package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	//load env
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Warning: .env file not found, using default environment variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//init router
	router := chi.NewRouter()
	//essential middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	//Health Check Routes
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			return
		}
	})

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
