-- name: CreateCheckResult :one
INSERT INTO check_results (monitor_id, status, status_code, latency_ms, error_msg)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCheckResultsByMonitorID :many
SELECT * FROM check_results
WHERE monitor_id = $1
ORDER BY checked_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUptimeStats :one
SELECT
    COUNT(*) AS total,
    COUNT(*) FILTER (WHERE status = 'up') AS total_up
FROM check_results
WHERE monitor_id = $1
  AND checked_at >= NOW() - INTERVAL '24 hours';