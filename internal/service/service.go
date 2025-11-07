package service

import (
	"context"

	"github.com/Vladimirmoscow84/Shortener.git/internal/model"
)

type shortenerRepo interface {
	SaveShortURL(ctx context.Context, short *model.ShortURL) (int, error)
	GetShortURL(ctx context.Context, shortCode string) (*model.ShortURL, error)
	SaveClick(ctx context.Context, click *model.Click) error
	GetClicksByShortURL(ctx context.Context, shortURLID int) ([]model.Click, error)
	GetClickAnalytics(ctx context.Context, shortURLID int) (map[string]map[string]int, error)
}

type cacheRepo interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type ServiceURL struct {
	storage shortenerRepo
	cache   cacheRepo
}

func New(storage shortenerRepo, cache cacheRepo) *ServiceURL {
	return &ServiceURL{
		storage: storage,
		cache:   cache,
	}
}
