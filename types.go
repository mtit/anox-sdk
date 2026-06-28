package sdk

import "time"

type LogLevel string

const (
	LogLevelDebug     LogLevel = "Debug"
	LogLevelInfo      LogLevel = "Info"
	LogLevelImportant LogLevel = "Important"
	LogLevelEmergency LogLevel = "Emergency"
)

type LogEntry struct {
	Time     time.Time         `json:"time"`
	Service  string            `json:"service"`
	Instance string            `json:"instance"`
	Level    LogLevel          `json:"level"`
	Action   string            `json:"action"`
	Message  string            `json:"message"`
	TraceID  string            `json:"trace_id,omitempty"`
	Stacks   []string          `json:"stacks,omitempty"`
	Context  map[string]string `json:"context,omitempty"`
}

type HeartbeatMessage struct {
	Type           string  `json:"type"`
	CPUCores       int     `json:"cpu_cores"`
	CPUPercent     float64 `json:"cpu_percent"`
	MemoryTotalMB  int64   `json:"memory_total_mb"`
	MemoryAvailMB  int64   `json:"memory_avail_mb"`
	GlobalVersion  int64   `json:"global_version"`
	ServiceVersion int64   `json:"service_version"`
}

type RegisterMessage struct {
	Type        string `json:"type"`
	ServiceName string `json:"service_name"`
	HttpHost    string `json:"http_host"`
	HttpPort    string `json:"http_port"`
}

type ConfigPullRequest struct {
	Type    string `json:"type"`
	Service string `json:"service"`
}
