package storage

import (
	"log"

	"github.com/Vladimirmoscow84/Shortener.git/internal/storage/cache"
	"github.com/jmoiron/sqlx"
	"github.com/wb-go/wbf/redis"
)

// Storage - структура для работы с БД и кэш
type Storage struct {
	DB    *sqlx.DB     //поле postgres
	Cache *cache.Cache //поле redis
}

// New - конструктор для создания экземпляра Storage
func New(databaseUri, rdAddr string) (*Storage, error) {
	db, err := sqlx.Connect("pgx", databaseUri)
	if err != nil {
		log.Fatalf("[storage] error connection to DB: %v", err)
	}
	rd := redis.New(rdAddr, "", 0)
	return &Storage{
		DB:    db,
		Cache: cache.NewCache(rd),
	}, nil
}

// написать отдельный пакет postgres с методами,как и для кэш
