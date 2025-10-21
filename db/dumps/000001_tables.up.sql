BEGIN;

CREATE TABLE IF NOT EXISTS short_url(
    id SERIAL KEY PRIMARY KEY,
    original_code TEXT NOT NULL,
    short_code VARCHAR(16) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS clics(
    id SERIAL PRIMARY KEY,
    short_url_id INTEGER NOT NULL REFERENCES short_urls(id) ON DELETE CASCADE,
    user_agent TEXT,
    temestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_short_urls_short_code ON short_urls(short_code);
CREATE INDEX idx_clicks_short_url_id ON clicks(short_url_id);

COMMIT;