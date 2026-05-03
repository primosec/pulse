package domain

import "time"

type CheckStatus string

const (
	CheckStatusUp      CheckStatus = "up"
	CheckStatusDown    CheckStatus = "down"
	CheckStatusTimeout CheckStatus = "timeout"
)

type CheckResult struct {
	ID         string      `json:"id"`
	MonitorID  string      `json:"monitor_id"`
	Status     CheckStatus `json:"status"`
	StatusCode *int        `json:"status_code"`
	Latency    *int        `json:"latency"`
	ErrorMsg   *string     `json:"error_msg"`
	CheckedAt  time.Time   `json:"checked_at"`
}
