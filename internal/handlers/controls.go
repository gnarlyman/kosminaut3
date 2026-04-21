package handlers

import (
	"net/http"
	"strconv"

	"kosminaut3/internal/iss"
	"kosminaut3/internal/views"
)

func Controls(r *views.Renderer, client *iss.Client, defaultPollSec int) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := req.ParseForm(); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}

		interval := parseIntervalForm(req, defaultPollSec)
		paused := req.FormValue("paused") == "true"

		var (
			pos     iss.Position
			fetchErr error
		)
		if !paused {
			pos, fetchErr = client.Fetch(req.Context())
		}
		view := newIssView(pos, interval, paused, fetchErr)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if rErr := r.Partial(w, "iss", view); rErr != nil {
			http.Error(w, "render error", http.StatusInternalServerError)
		}
	}
}

func parseIntervalForm(req *http.Request, fallback int) int {
	if raw := req.FormValue("interval"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 && n <= 60 {
			return n
		}
	}
	return fallback
}
