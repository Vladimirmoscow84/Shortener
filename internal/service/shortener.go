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

// ReturnOriginalURLByShort возвращает оригинальную ссылку по короткой
func (s *ServiceURL) ReturnOriginalURLByShort(ctx context.Context, shortCode, userAgent string) (string, error) {

	var short *model.ShortURL

	longUrl, err := s.cache.Get(ctx, shortCode)
	if err == nil && longUrl != "" {
		log.Printf("[service] availible in cache %s", shortCode)
	} else {
		log.Println("[service] failed to get short URL by cache")
		log.Println("[service] traying to get short URL from DB")

		short, err := s.storage.GetShortURL(ctx, shortCode)
		if err != nil {
			log.Println("[service] failed to get short URL by DB")
			return "", fmt.Errorf("[service] failed to get short URL by DB %w", err)
		}
		if short == nil {
			return "", fmt.Errorf("[service] short URL not found %s", shortCode)
		}
		longUrl = short.OriginalCode

		err = s.cache.Set(ctx, shortCode, longUrl)
		if err != nil {
			log.Printf("[service] failed to cache short code %s: %v", shortCode, err)

		}
	}

	click := &model.Click{
		ShortURLID: 0,
		UserAgent:  userAgent,
		Timestamp:  time.Now(),
	}

	if short != nil {
		click.ShortURLID = short.ID
	} else {
		dbShort, err := s.storage.GetShortURL(ctx, shortCode)
		if err != nil {
			log.Printf("[service] failed to get short url ID for click log: %v", err)
		} else if dbShort != nil {
			click.ShortURLID = dbShort.ID
		}
	}
	if click.ShortURLID != 0 {
		if err := s.storage.SaveClick(ctx, click); err != nil {
			log.Printf("[service] failed to log click for %s: %v", shortCode, err)
		}
	} else {
		log.Printf("[service] click not saved — no ShortURLID for %s", shortCode)
	}

	return longUrl, nil

}

// GetAnalytics возвращает статистику по кликам
func (s *ServiceURL) GetAnalytics(ctx context.Context, shortURLID uint) (map[string]map[string]int, error) {
	analytics, err := s.storage.GetClickAnalytics(ctx, shortURLID)
	if err != nil {
		return nil, fmt.Errorf("[service] failed to get analytics: %w", err)
	}
	return analytics, nil
}
