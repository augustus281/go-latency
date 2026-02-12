package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/augustus281/go-latency/singleflight/internal/metrics"
)

type TemplateService interface {
	GetTemplate(ctx context.Context, id string) (string, error)
}

// NewHandler is a function that creates a new HTTP handler for the template service
func NewHandler(s TemplateService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/template-details", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		metrics.InFlightRequests.Inc()
		defer metrics.InFlightRequests.Dec()

		status := "200"
		defer func() {
			duration := time.Since(start).Seconds()
			metrics.HTTPRequests.WithLabelValues("/template-details", r.Method, status).Inc()
			metrics.HTTPLatency.WithLabelValues("/template-details", r.Method).Observe(duration)
		}()

		id := r.URL.Query().Get("id")
		if id == "" {
			status = "400"
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		result, err := s.GetTemplate(r.Context(), id)
		if err != nil {
			status = "500"
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(result))
	})

	return mux
}
