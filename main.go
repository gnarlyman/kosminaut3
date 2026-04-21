package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"kosminaut3/internal/config"
	"kosminaut3/internal/iss"
	"kosminaut3/internal/server"
	"kosminaut3/internal/views"
)

func main() {
	cfg := config.Load()

	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatalf("sub fs: %v", err)
	}
	staticFS, err := fs.Sub(sub, "static")
	if err != nil {
		log.Fatalf("sub static: %v", err)
	}

	renderer, err := views.New(sub)
	if err != nil {
		log.Fatalf("views: %v", err)
	}

	srv := server.New(server.Deps{
		Cfg:      cfg,
		Renderer: renderer,
		Client:   iss.NewClient(),
		StaticFS: staticFS,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("kosminaut3 listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
