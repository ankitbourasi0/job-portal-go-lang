package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ankitbourasi0/job-portal/internal/database"
	"github.com/ankitbourasi0/job-portal/internal/handler"
	"github.com/ankitbourasi0/job-portal/internal/repository"
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

	//load url from .env
	dbURL := os.Getenv("SUPABASE_LOCAL_DB_URL")
	if dbURL == "" {
		log.Fatal("SUPABASE_LOCAL_DB_URL environment variable not set")
	}

	//DB connection establish
	db := database.InitDB(dbURL)
	defer db.Close() //close the database connection when server shutdown

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// SQLC Queries initialize karein
	queries := database.New(db)

	//Initialize Repository & Handler
	jobRepo := repository.NewJobRepository(queries)
	jobHandler := &handler.JobHandler{Repo: jobRepo}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// API Route
	router.Post("/api/jobs", jobHandler.HandleCreateJob)

	//Health Check Routes
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Database Connection SUCCESSFUL!"))
	})

	log.Printf("Server starting on port http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
