package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"subscriptions-app/internal/config"
	"subscriptions-app/internal/database"
	"subscriptions-app/internal/handlers"
	"subscriptions-app/internal/middleware"
	"subscriptions-app/internal/repository"
	"subscriptions-app/internal/validator"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db, "database/schema.sql"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	validator := validator.New()
	userRepo := repository.NewUserRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)

	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(userRepo, validator, cfg.JWTSecret, 24*time.Hour)
	subHandler := handlers.NewSubscriptionHandler(subRepo, validator)

	// === СОЗДАЁМ РОУТЕР ===
	r := chi.NewRouter()

	// Middleware от chi
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", healthHandler.Check)
	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		r.Get("/me", authHandler.Me)

		r.Route("/subscriptions", func(r chi.Router) {
			r.Get("/", subHandler.GetAll)
			r.Post("/", subHandler.Create)
			r.Get("/{id}", subHandler.GetByID)
			r.Put("/{id}", subHandler.Update)
			r.Delete("/{id}", subHandler.Delete)
		})

		r.Get("/stats", subHandler.GetStats)
	})

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{
			"*",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            true,
	})

	// === 2. ОБОРАЧИВАЕМ РОУТЕР ===
	handler := corsMiddleware.Handler(r)

	log.Printf("Server starting on %s", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, handler); err != nil { // ← handler, НЕ r!
		log.Fatalf("Server failed: %v", err)
	}
}
