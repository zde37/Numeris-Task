package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/zde37/Numeris-Task/internal/config"
	"github.com/zde37/Numeris-Task/internal/controller"
	"github.com/zde37/Numeris-Task/internal/repository"
	"github.com/zde37/Numeris-Task/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run is the main entry point for the application. It sets up the application configuration,
// initializes the database connection, creates the service and handler instances, starts the
// HTTP server, and handles the graceful shutdown of the server.
func run() error {
	cfg := config.Load(os.Getenv("ENVIRONMENT"), os.Getenv("HTTP_SERVER_ADDRESS"),
		os.Getenv("DSN"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := config.SetupDatabase(ctx, cfg, "file://migrations")
	if err != nil {
		return err
	}
	defer dbPool.Close()

	repo := repository.NewRepository(dbPool)
	srvc := service.NewService(repo)
	hndl := controller.NewHandlerImpl(cfg.Environment, srvc)

	srv := &http.Server{
		Addr:    cfg.HTTPServerAddr,
		Handler: hndl.GetRouter(),
	}

	go func() {
		log.Printf("server started on %s", cfg.HTTPServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to start server: %v", err)
			cancel()
		}
	}()

	return gracefulShutdown(ctx, srv)
}

// gracefulShutdown is a function that handles the graceful shutdown of an HTTP server. 
func gracefulShutdown(ctx context.Context, srv *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("server gracefully stopped")
	return nil
}
