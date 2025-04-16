CREATE TABLE IF NOT EXISTS refreshSessions (
    id SERIAL PRIMARY KEY,
    userId UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    refreshToken VARCHAR(200) NOT NULL,
    expiresAt TIMESTAMP WITH TIME ZONE NOT NULL,
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    ip VARCHAR(15)
)

CREATE INDEX idx_refresh_tokens_userId ON refresh_tokens(userId);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(refreshToken);