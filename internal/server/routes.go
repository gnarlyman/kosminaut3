package server

import (
	"net/http"

	"kosminaut3/internal/handlers"
)

func Register(mux *http.ServeMux, deps Deps) {
	poll := deps.Cfg.DefaultPollSec

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(deps.StaticFS))))
	mux.Handle("GET /iss", handlers.ISS(deps.Renderer, deps.Client, poll))
	mux.Handle("POST /controls", handlers.Controls(deps.Renderer, deps.Client, poll))
	mux.Handle("GET /", handlers.Index(deps.Renderer, poll))
}
