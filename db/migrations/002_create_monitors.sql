CREATE TABLE monitors (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    url          TEXT NOT NULL,
    type         TEXT NOT NULL DEFAULT 'http',
    interval     INT NOT NULL DEFAULT 60,
    timeout      INT NOT NULL DEFAULT 10,
    status       TEXT NOT NULL DEFAULT 'pending',
    uptime_pct   FLOAT NOT NULL DEFAULT 0,
    last_checked TIMESTAMPTZ,
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_monitors_user_id ON monitors(user_id);
CREATE INDEX idx_monitors_is_active ON monitors(is_active);