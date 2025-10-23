package storage

import (
	"errors"
	"log"

	"github.com/Vladimirmoscow84/Shortener.git/internal/storage/postgres"
	"github.com/wb-go/wbf/redis"
)

// Storage - структура для работы с БД и кэш
type Storage struct {
	*postgres.Postgres
	*redis.Client
}

// New - конструктор для создания экземпляра Storage
func New(pg *postgres.Postgres, rd *redis.Client) (*Storage, error) {
	if pg == nil {
		log.Println("[storage] postgres client is nil ")
		return nil, errors.New("[storage] postgres client is nill")
	}
	if rd == nil {
		log.Println("[storage] redis client is nil ")
		return nil, errors.New("[storage] redis client is nill")
	}
	return &Storage{
		Postgres: pg,
		Client:   rd,
	}, nil
}
