package handlers

import (
	"net/http"
	"strconv"

	"kosminaut3/internal/iss"
	"kosminaut3/internal/views"
)

func ISS(r *views.Renderer, client *iss.Client, defaultPollSec int) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		interval := parseIntervalQuery(req, defaultPollSec)
		paused := req.URL.Query().Get("paused") == "true"

		pos, err := client.Fetch(req.Context())
		view := newIssView(pos, interval, paused, err)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if rErr := r.Partial(w, "iss", view); rErr != nil {
			http.Error(w, "render error", http.StatusInternalServerError)
		}
	}
}

func parseIntervalQuery(req *http.Request, fallback int) int {
	if raw := req.URL.Query().Get("interval"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 && n <= 60 {
			return n
		}
	}
	return fallback
}
