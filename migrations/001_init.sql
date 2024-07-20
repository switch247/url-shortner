CREATE TABLE urls (
    id TEXT PRIMARY KEY,
    long_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    click_count BIGINT DEFAULT 0
);
-- no need for this migrations if gorm  auto migration is used

CREATE TABLE clicks (
    id BIGSERIAL PRIMARY KEY,
    short_code TEXT NOT NULL,
    ip TEXT,
    country TEXT,
    city TEXT,
    user_agent TEXT,
    referrer TEXT,
    device TEXT,
    os TEXT,
    browser TEXT,
    utms TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_clicks_code ON clicks(short_code);
CREATE INDEX idx_clicks_created ON clicks(created_at);