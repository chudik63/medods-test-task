CREATE TABLE IF NOT EXISTS refreshSessions (
    id SERIAL PRIMARY KEY,
    userId UUID REFERENCES users(id) ON DELETE CASCADE,
    refreshToken VARCHAR(200) NOT NULL,
    ip VARCHAR(15) NOT NULL,
    expiresIn bigint NOT NULL,
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX idx_refresh_tokens_userId ON refresh_tokens(userId);