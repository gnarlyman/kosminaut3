package server

import (
	"io/fs"
	"net/http"
	"time"

	"kosminaut3/internal/config"
	"kosminaut3/internal/iss"
	"kosminaut3/internal/views"
)

type Deps struct {
	Cfg      config.Config
	Renderer *views.Renderer
	Client   *iss.Client
	StaticFS fs.FS
}

func New(deps Deps) *http.Server {
	mux := http.NewServeMux()
	Register(mux, deps)

	return &http.Server{
		Addr:              ":" + deps.Cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
