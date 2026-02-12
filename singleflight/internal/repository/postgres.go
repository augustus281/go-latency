package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TemplateRepository struct {
	db *pgxpool.Pool
}

func NewTemplateRepository(db *pgxpool.Pool) *TemplateRepository {
	return &TemplateRepository{
		db: db,
	}
}

// GetByID is a method that returns the name of a template by its ID
func (r *TemplateRepository) GetByID(ctx context.Context, id string) (string, error) {
	var name string
	if err := r.db.QueryRow(ctx, "SELECT name FROM templates WHERE id = $1", id).Scan(&name); err != nil {
		slog.Error("error to found template's name", "error", err)
		return "", fmt.Errorf("failed to get template's name by id, %w", err)
	}
	return name, nil
}
