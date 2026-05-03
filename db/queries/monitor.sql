-- name: CreateMonitor :one
INSERT INTO monitors (user_id, name, url, type, interval, timeout)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetMonitorByID :one
SELECT * FROM monitors
WHERE id = $1
LIMIT 1;

-- name: GetMonitorsByUserID :many
SELECT * FROM monitors
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetActiveMonitors :many
SELECT * FROM monitors
WHERE is_active = TRUE
ORDER BY id;

-- name: UpdateMonitorStatus :one
UPDATE monitors
SET status = $2, uptime_pct = $3, last_checked = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateMonitor :one
UPDATE monitors
SET name = $2, url = $3, interval = $4, timeout = $5
WHERE id = $1 AND user_id = $6
RETURNING *;

-- name: ToggleMonitorActive :one
UPDATE monitors
SET is_active = $2
WHERE id = $1 AND user_id = $3
RETURNING *;

-- name: DeleteMonitor :exec
DELETE FROM monitors
WHERE id = $1 AND user_id = $2;