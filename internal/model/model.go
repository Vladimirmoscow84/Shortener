package model

import "time"

type ShortURL struct {
	ID           uint
	OriginalCode string
	ShortCode    string
	CreatedAt    time.Time
	Clicks       []Click
}

type Click struct {
	ID         uint
	ShortURLID uint
	UserAgent  string
	Timestamp  time.Time
}
