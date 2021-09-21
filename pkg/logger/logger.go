package logger

import (
	"encoding/json"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Levels
type level string

//
const (
	levelInfo  level = "INFO"
	levelError level = "ERROR"
	levelFetal level = "FETAL"
)

// Logger
type Logger struct {
	file *os.File
	mu   *sync.Mutex
}

// NewLogger
func NewLogger(path string) (*Logger, error) {
	// open or create log
	log, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{
		file: log,
		mu:   &sync.Mutex{},
	}, nil
}

// print
func (l *Logger) print(level level, on, message string, properties map[string]string) (int, error) {
	// If the severity level of the log entry is below the minimum severity for the
	// log row
	row := struct {
		Time       string            `json:"time"`
		Level      string            `json:"level"`
		Function   string            `json:"function"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Time:       time.Now().UTC().Format(time.RFC3339),
		Level:      string(level),
		Function:   on,
		Message:    message,
		Properties: properties,
	}

	// stack trace
	if level == levelFetal {
		row.Trace = string(debug.Stack())
	}

	// Declare a line variable for holding the actual log entry text.
	var line []byte
	// Marshal the anonymous struct to JSON and store it in the line variable. If there // was a problem creating the JSON, set the contents of the log entry to be that
	// plain-text error message instead.
	line, err := json.Marshal(row)
	if err != nil {
		line = []byte(string(levelError) + ": unable to marshal log message:" + err.Error())
	}
	// Lock the mutex so that no two writes to the output destination cannot happen // concurrently. If we don't do this, it's possible that the text for two or more // log entries will be intermingled in the output.
	l.mu.Lock()
	defer l.mu.Unlock()

	// Write the log entry followed by a newline.
	return l.file.Write(append(line, '\n'))
}

// PrintInfo
func (l *Logger) Info(on, message string, properties map[string]string) {
	l.print(levelInfo, on, message, properties)
}

// PrintError
func (l *Logger) Error(on, message string, properties map[string]string) {
	l.print(levelError, on, message, properties)
}

// Fetal
func (l *Logger) Fetal(on, message string, properties map[string]string) {
	l.print(levelFetal, on, message, properties)
}

// Close
func (l *Logger) Close() {
	l.file.Close()
}
