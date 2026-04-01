package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/internal/api/v1"
	"go-reading-log-api-next/internal/api/v1/middleware"
	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultShutdownTimeout = 5 * time.Second

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	log := logger.Initialize(cfg.LogLevel, cfg.LogFormat)

	// Log startup information
	log.Info("Starting server...",
		"host", cfg.ServerHost,
		"port", cfg.ServerPort,
		"log_level", cfg.LogLevel,
		"log_format", cfg.LogFormat)

	// Build database connection string
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabase)

	// Connect to database with connection pooling
	log.Info("Connecting to database...")
	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Verify database connection
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Error("Database connection ping failed", "error", err)
		os.Exit(1)
	}
	log.Info("Database connection established")

	// Create repository instances with the connection pool
	projectRepo := postgres.NewProjectRepositoryImpl(dbPool)
	logRepo := postgres.NewLogRepositoryImpl(dbPool)

	// Setup routes with repositories
	router := api.SetupRoutes(projectRepo, logRepo)

	// Create middleware chain: Recovery -> CORS -> RequestID -> Logging -> Handler
	middlewareChain := middleware.Chain(router,
		middleware.RecoveryMiddleware,
		middleware.CORSMiddleware,
		middleware.RequestIDMiddleware,
		middleware.LoggingMiddleware,
	)

	// Create HTTP server with timeout settings
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
		Handler:      middlewareChain,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info("Server starting", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	log.Info("Server is ready to accept connections")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Block until signal received
	<-quit
	log.Info("Shutdown signal received")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer shutdownCancel()

	// Graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown error", "error", err)
		// Force close if shutdown times out
		if err := server.Close(); err != nil {
			log.Error("Server close error", "error", err)
		}
	}

	log.Info("Server stopped gracefully")
}
