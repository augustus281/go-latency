package handler

import (
	"context"
	"net/http"
)

type TemplateService interface {
	GetTemplate(ctx context.Context, id string) (string, error)
}

// NewHandler is a function that creates a new HTTP handler for the template service
func NewHandler(s TemplateService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/template-details", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		result, err := s.GetTemplate(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(result))
	})

	return mux
}
