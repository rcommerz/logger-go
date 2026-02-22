package logger

import "time"

// LogLevel represents the severity of a log entry
type LogLevel string

const (
	LevelINFO  LogLevel = "INFO"
	LevelERROR LogLevel = "ERROR"
	LevelWARN  LogLevel = "WARN"
	LevelDEBUG LogLevel = "DEBUG"
)

// LogType represents the category of a log entry
type LogType string

const (
	TypeNormal   LogType = "normal"
	TypeHTTP     LogType = "http"
	TypeError    LogType = "error"
	TypeSecurity LogType = "security"
	TypeAudit    LogType = "audit"
	TypeDebug    LogType = "debug"
)

// Config holds logger initialization configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Env            string
	Level          LogLevel
}

// LogContext holds arbitrary key-value pairs for structured logging
type LogContext map[string]interface{}

// Fields is a helper function to create LogContext from alternating key-value pairs
// Example: Fields("key1", "value1", "key2", "value2")
func Fields(keysAndValues ...interface{}) LogContext {
	if len(keysAndValues)%2 != 0 {
		panic("Fields requires an even number of arguments")
	}

	context := make(LogContext)
	for i := 0; i < len(keysAndValues)-1; i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			panic("Field keys must be strings")
		}
		context[key] = keysAndValues[i+1]
	}
	return context
}

// MeasureDuration calculates the duration in milliseconds since the given start time
func MeasureDuration(start time.Time) float64 {
	return float64(time.Since(start).Milliseconds())
}
