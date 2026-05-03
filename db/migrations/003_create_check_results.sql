CREATE TABLE check_results (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    monitor_id  UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    status      TEXT NOT NULL,
    status_code INT,
    latency_ms  INT,
    error_msg   TEXT,
    checked_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_check_results_monitor_id ON check_results(monitor_id);
CREATE INDEX idx_check_results_checked_at ON check_results(checked_at);