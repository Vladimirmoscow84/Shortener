package postgres_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Vladimirmoscow84/Shortener.git/internal/model"
	"github.com/Vladimirmoscow84/Shortener.git/internal/storage/postgres"
	"github.com/stretchr/testify/require"
)

//docker run -d --name pg-shortener -p 5440:5432 -e POSTGRES_USER=vladimir -e POSTGRES_PASSWORD="password" -e POSTGRES_DB="shortener" postgres:latest

const testDSN = "host=localhost port=5440 user=vladimir password=password dbname=shortener sslmode=disable"

func setupTestDB(t *testing.T) *postgres.Postgres {
	pg, err := postgres.New(testDSN)
	require.NoError(t, err)

	ctx := context.Background()

	_, _ = pg.DB.ExecContext(ctx, `DROP TABLE IF EXISTS clicks;`)
	_, _ = pg.DB.ExecContext(ctx, `DROP TABLE IF EXISTS short_urls;`)

	schema := `
	CREATE TABLE short_urls (
		id SERIAL PRIMARY KEY,
		original_code TEXT NOT NULL,
		short_code TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL
	);

	CREATE TABLE clicks (
		id SERIAL PRIMARY KEY,
		short_url_id INTEGER REFERENCES short_urls(id) ON DELETE CASCADE,
		user_agent TEXT NOT NULL,
		timestamp TIMESTAMP NOT NULL
	);
	`
	_, err = pg.DB.ExecContext(ctx, schema)
	require.NoError(t, err)

	return pg
}

func TestPostgresStorage(t *testing.T) {
	pg := setupTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	//SaveShortURL
	short := &model.ShortURL{
		OriginalCode: "https://example.com/long-url",
		ShortCode:    "abc123",
		CreatedAt:    time.Now(),
	}
	id, err := pg.SaveShortURL(ctx, short)
	require.NoError(t, err)
	require.Greater(t, id, 0, "ID должен быть больше 0")

	//GetShortURL
	found, err := pg.GetShortURL(ctx, short.ShortCode)
	require.NoError(t, err)
	require.NotNil(t, found)
	require.Equal(t, short.OriginalCode, found.OriginalCode)
	require.Equal(t, short.ShortCode, found.ShortCode)

	//SaveClick
	click := &model.Click{
		ShortURLID: found.ID,
		UserAgent:  "Mozilla/5.0 (Windows NT 10.0)",
		Timestamp:  time.Now(),
	}
	err = pg.SaveClick(ctx, click)
	require.NoError(t, err)
	require.Greater(t, click.ID, 0)

	//GetClicksByShortURL
	clicks, err := pg.GetClicksByShortURL(ctx, found.ID)
	require.NoError(t, err)
	require.Len(t, clicks, 1)
	require.Equal(t, click.UserAgent, clicks[0].UserAgent)

	//GetClickAnalytics
	analytics, err := pg.GetClickAnalytics(ctx, found.ID)
	require.NoError(t, err)
	require.NotEmpty(t, analytics, "аналитика не должна быть пустой")

	for day, uaMap := range analytics {
		log.Printf("День: %s", day)
		for ua, count := range uaMap {
			log.Printf("UA: %s — %d кликов", ua, count)
		}
	}
}
