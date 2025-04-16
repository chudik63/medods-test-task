CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(100)
);


CREATE TABLE IF NOT EXISTS refreshSessions (
    id SERIAL PRIMARY KEY,
    userId UUID REFERENCES users(id) ON DELETE CASCADE,
    ip VARCHAR(45),
    refreshToken VARCHAR(200) NOT NULL,
    expiresAt TIMESTAMP WITH TIME ZONE NOT NULL,
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX idx_refresh_tokens_userId ON refreshSessions(userId);
CREATE INDEX idx_refresh_tokens_token ON refreshSessions(refreshToken);