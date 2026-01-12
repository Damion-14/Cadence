package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/damion-14/cadence/backend/internal/cache"
	"github.com/damion-14/cadence/backend/internal/config"
	"github.com/damion-14/cadence/backend/internal/database"
	"github.com/damion-14/cadence/backend/internal/handlers"
	"github.com/damion-14/cadence/backend/internal/middleware"
	"github.com/damion-14/cadence/backend/internal/router"
	"github.com/damion-14/cadence/backend/internal/services"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	db, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()
	fmt.Println("Connected to PostgreSQL")

	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	defer redisClient.Close()
	fmt.Println("Connected to Redis")

	cacheClient := cache.NewCache(redisClient)

	workoutService := services.NewWorkoutService(db, cacheClient)
	statsService := services.NewStatsService(db, cacheClient)

	deps := &router.Dependencies{
		DB:              db,
		Redis:           redisClient,
		Cache:           cacheClient,
		Config:          cfg,
		AuthHandler:     handlers.NewAuthHandler(db, cfg.JWT),
		WorkoutHandler:  handlers.NewWorkoutHandler(workoutService),
		ExerciseHandler: handlers.NewExerciseHandler(workoutService),
		StatsHandler:    handlers.NewStatsHandler(statsService),
	}

	mux := router.NewRouter(deps)

	handler := middleware.Recovery(
		middleware.Logger(
			middleware.RequestID(
				middleware.CORS(cfg.CORS.AllowedOrigins)(mux),
			),
		),
	)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdownChan := make(chan error, 1)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		fmt.Println("\nShutting down server...")
		shutdownChan <- server.Close()
	}()

	fmt.Printf("Server listening on port %s\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return <-shutdownChan
}
