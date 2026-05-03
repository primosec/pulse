package domain

import "time"

type MonitorType string
type MonitorStatus string

const (
	MonitorTypeHTTP MonitorType = "http"
	MonitorTypeTCP  MonitorType = "tcp"
	MonitorTypePing MonitorType = "ping"
)

const (
	MonitorStatusPending MonitorStatus = "pending"
	MonitorStatusUp      MonitorStatus = "up"
	MonitorStatusDown    MonitorStatus = "down"
)

type Monitor struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	Name        string        `json:"name"`
	URL         string        `json:"url"`
	Type        MonitorType   `json:"type"`
	Interval    int           `json:"interval"`
	Timeout     int           `json:"timeout"`
	Status      MonitorStatus `json:"status"`
	UptimePct   float64       `json:"uptime_pct"`
	LastChecked *time.Time    `json:"last_checked"`
	IsActive    bool          `json:"is_active"`
	CreatedAt   time.Time     `json:"created_at"`
}
