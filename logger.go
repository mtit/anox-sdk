package sdk

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

type Logger struct {
	client    *Client
	buffer    []LogEntry
	mu        sync.Mutex
	batchSize int
	timeout   time.Duration
	stopCh    chan struct{}
}

func newLogger(client *Client) *Logger {
	l := &Logger{client: client, buffer: make([]LogEntry, 0, 10), batchSize: 10, timeout: 20 * time.Second, stopCh: make(chan struct{})}
	go l.flushLoop()
	return l
}

func (l *Logger) Log(entry LogEntry, display bool) {
	if display {
		l.displayLog(entry)
	}
	l.mu.Lock()
	l.buffer = append(l.buffer, entry)
	shouldFlush := len(l.buffer) >= l.batchSize
	l.mu.Unlock()
	if shouldFlush {
		l.Flush()
	}
}

func (l *Logger) LogDebug(action string, message string, ctx map[string]string, display bool) {
	l.Log(LogEntry{Time: time.Now(), Service: l.client.GetServiceName(), Instance: l.client.GetInstanceID(), Level: LogLevelDebug, Action: action, Message: message, Context: ctx}, display)
}

func (l *Logger) LogInfo(action string, message string, ctx map[string]string, display bool) {
	l.Log(LogEntry{Time: time.Now(), Service: l.client.GetServiceName(), Instance: l.client.GetInstanceID(), Level: LogLevelInfo, Action: action, Message: message, Context: ctx}, display)
}

func (l *Logger) LogImportant(action string, message string, ctx map[string]string, display bool) {
	l.Log(LogEntry{Time: time.Now(), Service: l.client.GetServiceName(), Instance: l.client.GetInstanceID(), Level: LogLevelImportant, Action: action, Message: message, Context: ctx}, display)
}

func (l *Logger) LogEmergency(action string, message string, ctx map[string]string, display bool) {
	l.Log(LogEntry{Time: time.Now(), Service: l.client.GetServiceName(), Instance: l.client.GetInstanceID(), Level: LogLevelEmergency, Action: action, Message: message, Stacks: []string{string(debug.Stack())}, Context: ctx}, display)
}

func (l *Logger) LogError(action string, err error, ctx map[string]string, display bool) {
	l.Log(LogEntry{Time: time.Now(), Service: l.client.GetServiceName(), Instance: l.client.GetInstanceID(), Level: LogLevelEmergency, Action: action, Message: err.Error(), Stacks: []string{string(debug.Stack())}, Context: ctx}, display)
}

func (l *Logger) Flush() {
	l.mu.Lock()
	if len(l.buffer) == 0 {
		l.mu.Unlock()
		return
	}
	logs := make([]LogEntry, len(l.buffer))
	copy(logs, l.buffer)
	l.buffer = l.buffer[:0]
	l.mu.Unlock()
	l.sendLogs(logs)
}

func (l *Logger) Close() {
	close(l.stopCh)
	l.Flush()
}

func (l *Logger) flushLoop() {
	ticker := time.NewTicker(l.timeout)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			l.Flush()
		case <-l.stopCh:
			return
		}
	}
}

func (l *Logger) sendLogs(logs []LogEntry) {
	if len(logs) == 0 {
		return
	}
	instanceID := l.client.GetInstanceID()
	logsData := make([]map[string]interface{}, len(logs))
	for i, logEntry := range logs {
		logsData[i] = map[string]interface{}{
			"time":     logEntry.Time.Format(time.RFC3339Nano),
			"service":  logEntry.Service,
			"instance": instanceID,
			"level":    logEntry.Level,
			"action":   logEntry.Action,
			"message":  logEntry.Message,
			"trace_id": logEntry.TraceID,
			"stacks":   logEntry.Stacks,
			"context":  logEntry.Context,
		}
	}
	msg := map[string]interface{}{"type": "logs_batch", "logs": logsData}
	if err := l.client.writeJSON(msg); err != nil {
		log.Printf("[Anox SDK] Failed to send logs: %v", err)
		l.mu.Lock()
		l.buffer = append(logs, l.buffer...)
		l.mu.Unlock()
		go l.client.reconnect()
	}
}

func (l *Logger) displayLog(entry LogEntry) {
	var levelColor string
	switch entry.Level {
	case LogLevelDebug:
		levelColor = "\033[36m"
	case LogLevelInfo:
		levelColor = "\033[32m"
	case LogLevelImportant:
		levelColor = "\033[33m"
	case LogLevelEmergency:
		levelColor = "\033[31m"
	default:
		levelColor = "\033[0m"
	}
	resetColor := "\033[0m"
	fmt.Printf("%s[%s]%s %s[%s]%s %s - %s\n", levelColor, entry.Level, resetColor, levelColor, entry.Action, resetColor, entry.Time.Format("2006-01-02 15:04:05"), entry.Message)
	if len(entry.Stacks) > 0 {
		for _, stack := range entry.Stacks {
			fmt.Printf("%sStack:%s\n%s\n", levelColor, resetColor, stack)
		}
	}
}
