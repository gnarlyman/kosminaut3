package handlers

import (
	"net/http"

	"kosminaut3/internal/views"
)

type indexData struct {
	IntervalSec int
	Paused      bool
}

func Index(r *views.Renderer, defaultPollSec int) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		data := indexData{IntervalSec: defaultPollSec, Paused: false}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := r.Page(w, data); err != nil {
			http.Error(w, "render error", http.StatusInternalServerError)
		}
	}
}
