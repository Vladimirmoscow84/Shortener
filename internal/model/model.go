package model

import "time"

type ShortURL struct {
	ID           int       `json:"id" db:"id"`
	OriginalCode string    `json:"original_code" db:"original_code"`
	ShortCode    string    `json:"short_code" db:"short_code"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	Clicks       []Click   `json:"clicks"`
}

type Click struct {
	ID         int       `json:"id" db:"id"`
	ShortURLID int       `json:"short-url_id" db:"short_url_id"`
	UserAgent  string    `json:"user-agent" db:"user_agent"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
}
