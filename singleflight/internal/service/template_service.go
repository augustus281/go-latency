package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/augustus281/go-latency/singleflight/internal/metrics"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
}

type Repository interface {
	GetByID(ctx context.Context, id string) (string, error)
}

type TemplateService struct {
	cache Cache
	repo  Repository
	group singleflight.Group
}

func NewTemplateService(c Cache, r Repository) *TemplateService {
	return &TemplateService{
		cache: c,
		repo:  r,
	}
}

// GetTemplate is a method that gets a template by its ID
func (s *TemplateService) GetTemplate(ctx context.Context, id string) (string, error) {
	key := fmt.Sprintf("template:%s", id)

	val, err := s.cache.Get(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("failed to get template from cache, %w", err)
	}
	if err == nil {
		slog.Info("CACHE HIT for template", "id", id)
		metrics.CacheHits.Inc()
		return val, nil
	}

	metrics.CacheMisses.Inc()

	result, err, shared := s.group.Do(key, func() (interface{}, error) {
		slog.Info("DB HIT for template", "id", id)
		metrics.DBHits.Inc()

		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		val, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return "", fmt.Errorf("failed to get template from database, %w", err)
		}

		_ = s.cache.Set(ctx, key, val, 30*time.Second)
		return val, nil
	})

	if shared {
		slog.Info("Result is shared from cache")
		metrics.SingleflightShared.Inc()
	}

	if err != nil {
		return "", fmt.Errorf("failed to get template, %w", err)
	}

	return result.(string), nil
}
