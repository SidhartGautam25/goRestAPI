package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SidhartGautam25/goRestAPI/internal/config"
	"github.com/SidhartGautam25/goRestAPI/internal/http/handlers/student"
	"github.com/SidhartGautam25/goRestAPI/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// databse setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialized ", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGABRT)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("error,server not started")

		}

	}()
	<-done

	slog.Info("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
