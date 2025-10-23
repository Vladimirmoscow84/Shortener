package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Vladimirmoscow84/Shortener.git/internal/model"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	DB *sqlx.DB
}

func New(databaseURI string) (*Postgres, error) {
	db, err := sqlx.Connect("pgx", databaseURI)
	if err != nil {
		log.Fatalf("[postgres] failed connect to DB: %v", err)
		return nil, fmt.Errorf("failed to connect to DB %w", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("[postgres] ping failed: %v", err)
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	log.Println("[postgres] connected successfully")
	return &Postgres{
		DB: db,
	}, nil
}

// Close закрывает соединение с БД
func (p *Postgres) Close() error {
	if p.DB != nil {
		log.Println("[postgres] closing connection")
		return p.DB.Close()
	}
	return nil
}

// SaveShortURL сохраняет короткую и оригинальную ссылку в БД
func (p *Postgres) SaveShortURL(ctx context.Context, short *model.ShortURL) (int, error) {
	row := p.DB.QueryRowContext(ctx, `
		INSERT INTO short_urls
			(original_code, short_code, created_at)
		VALUES
			($1,$2,$3)
		RETURNING id;
	`, short.OriginalCode, short.ShortCode, short.CreatedAt)

	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Printf("[postgres] error adding shortURL to base: %v", err)
		return 0, fmt.Errorf("error adding to base %w", err)
	}
	short.ID = id
	return id, nil
}

// GetShortURL возвращает запись из БД  short_urls по short_code
func (p *Postgres) GetShortURL(ctx context.Context, shortCode string) (*model.ShortURL, error) {
	var short model.ShortURL
	err := p.DB.GetContext(ctx, &short, `
		SELECT id, original_code, short_code, created_at
		FROM short_urls
		WHERE short_code = $1;
	`, shortCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("[postgres] failed to get short_url %s from DB: %v", shortCode, err)
		return nil, fmt.Errorf("failed to get short_url %s from DB: %w", shortCode, err)
	}
	return &short, nil
}

// SaveClick сохраняет переход по короткой ссылке
func (p *Postgres) SaveClick(ctx context.Context, click *model.Click) error {
	row := p.DB.QueryRowContext(ctx, `
		INSERT INTO clicks
			(short_url_id, user_agent, timestamp)
		VALUES
			($1,$2,$3)
		RETURNING id;
	`, click.ShortURLID, click.UserAgent, click.Timestamp)

	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Printf("[postgres] failed to save click for short_url_id=%d: %v", click.ShortURLID, err)
		return fmt.Errorf("failed to save click: %w", err)
	}
	click.ID = id
	return nil
}

// GetClicksByShortURL возвращает все клики по короткой ссылке
func (p *Postgres) GetClicksByShortURL(ctx context.Context, shortURLID int) ([]model.Click, error) {
	var clicks []model.Click
	err := p.DB.SelectContext(ctx, &clicks, `
		SELECT id, short_url_id, user_agent, timestamp
		FROM clicks
		WHERE short_url_id=$1
		ORDER BY timestamp DESC;
 `, shortURLID)
	if err != nil {
		log.Printf("[postgres] failed to get clicks for short_url_id=%d: %v", shortURLID, err)
		return nil, fmt.Errorf("failed to get clicks: %w", err)
	}
	return clicks, nil
}

// GetClickAnalytics возвращает количество кликов по дням и User-Agent
func (p *Postgres) GetClickAnalytics(ctx context.Context, shortURLID uint) (map[string]map[string]int, error) {
	const query = `
		SELECT DATE(timestamp) AS day, user_agent, COUNT(*) AS clicks
		FROM clicks
		WHERE short_url_id = $1
		GROUP BY day, user_agent
		ORDER BY day ASC;
	`

	rows, err := p.DB.QueryxContext(ctx, query, shortURLID)
	if err != nil {
		log.Printf("[postgres] failed to get click analytics for short_url_id=%d: %v", shortURLID, err)
		return nil, fmt.Errorf("failed to get analytics: %w", err)
	}
	defer rows.Close()

	analytics := make(map[string]map[string]int) // day - user_agent - count

	for rows.Next() {
		var day time.Time
		var userAgent string
		var count int
		if err := rows.Scan(&day, &userAgent, &count); err != nil {
			return nil, fmt.Errorf("failed to scan analytics row: %w", err)
		}

		dayStr := day.Format("2006-01-02")
		if _, ok := analytics[dayStr]; !ok {
			analytics[dayStr] = make(map[string]int)
		}
		analytics[dayStr][userAgent] = count
	}

	return analytics, nil
}
