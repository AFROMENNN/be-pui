package main

import (
	"be-pui/config"
	"be-pui/db"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	gin.SetMode(cfg.Server.Mode)

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(dbConn); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	srv := &http.Server{
		Addr: ":" + cfg.Server.Port,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.Server.BaseURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Received shutdown signal. Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server failed to shut down gracefully: %v", err)
	}

	log.Println("Server stopped.")
}
