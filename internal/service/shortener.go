package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Vladimirmoscow84/Shortener.git/internal/model"
)

// chars - набор допустимых символов для формирования короткой ссылки
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// lengthShortCode - длина короткой ссылки
const lengthSortCode = 8

// GenerateCode генерирует случайный код
func GenerateCode(length int) string {
	result := make([]byte, length)
	maximum := big.NewInt(int64(len(chars)))

	for i := 0; i < length; i++ {
		m, err := rand.Int(rand.Reader, maximum)
		if err != nil {
			result[i] = chars[int(m.Int64())%len(chars)]
			continue
		}
		result[i] = chars[m.Int64()]
	}
	return string(result)
}

// CreateShortURL созадет короткую ссылку
func (s *ServiceURL) CreateShortURL(ctx context.Context, longURL string) (*model.ShortURL, error) {
	cached, err := s.cache.Get(ctx, longURL)
	if err == nil && cached != "" {
		log.Printf("[] found cahsed short URLcode in cache for %s: %s", longURL, cached)
		short := &model.ShortURL{
			OriginalCode: longURL,
			ShortCode:    cached,
		}
		return short, nil
	}
	shortCode := GenerateCode(lengthSortCode)
	short := &model.ShortURL{
		OriginalCode: longURL,
		ShortCode:    shortCode,
		CreatedAt:    time.Now(),
	}
	// добавление в БД
	id, err := s.storage.SaveShortURL(ctx, short)
	if err != nil {
		log.Println("[srvice]failed to save short url in DB")
		return nil, fmt.Errorf("[service]failed to save short url in DB: %w", &err)
	}
	short.ID = id

	//кэширование
	err = s.cache.Set(ctx, shortCode, longURL)
	if err != nil {
		log.Println("[service] failed to cache short code")
		return nil, fmt.Errorf("[service] failed to cache short code %s:  %v", shortCode, err)

	}
	log.Printf("[service] created short URL %s: %s", longURL, shortCode)
	return short, nil
}
