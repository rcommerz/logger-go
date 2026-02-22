package logger

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// Example: Testing with observable logs instead of stdout capture
func TestLoggingWithObserver(t *testing.T) {
	// Reset instance for clean test
	instance = nil
	once = sync.Once{}

	// Create observable core to capture logs
	observedZapCore, observedLogs := observer.New(zapcore.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	// Initialize our logger with custom config
	initialized := Initialize(Config{
		ServiceName:    "observable-test",
		ServiceVersion: "1.0.0",
		Env:            "test",
		Level:          LevelDEBUG,
	})

	// Replace the zap instance with our observable one
	initialized.zap = observedLogger

	t.Run("should log info with verifiable fields", func(t *testing.T) {
		// Log something
		initialized.Info(context.Background(), "Test message", Fields(
			"user_id", "123",
			"action", "login",
		))

		// Get all logged entries
		allLogs := observedLogs.All()

		if len(allLogs) == 0 {
			t.Fatal("Expected at least one log entry")
		}

		// Get the last log entry
		logEntry := allLogs[len(allLogs)-1]

		// Verify message
		if logEntry.Message != "Test message" {
			t.Errorf("Expected message 'Test message', got '%s'", logEntry.Message)
		}

		// Verify level
		if logEntry.Level != zapcore.InfoLevel {
			t.Errorf("Expected INFO level, got %v", logEntry.Level)
		}

		// Verify fields
		fields := logEntry.Context
		foundUserId := false
		foundAction := false
		foundLogType := false

		for _, field := range fields {
			switch field.Key {
			case "user_id":
				foundUserId = true
				if field.String != "123" {
					t.Errorf("Expected user_id=123, got %s", field.String)
				}
			case "action":
				foundAction = true
				if field.String != "login" {
					t.Errorf("Expected action=login, got %s", field.String)
				}
			case "log_type":
				foundLogType = true
				if field.String != "normal" {
					t.Errorf("Expected log_type=normal, got %s", field.String)
				}
			}
		}

		if !foundUserId {
			t.Error("Expected user_id field")
		}
		if !foundAction {
			t.Error("Expected action field")
		}
		if !foundLogType {
			t.Error("Expected log_type field")
		}
	})

	t.Run("should log error with error type", func(t *testing.T) {
		// Clear previous logs
		observedLogs.TakeAll()

		// Log an error
		initialized.Error(context.Background(), "An error occurred", Fields(
			"code", "ERR_500",
			"endpoint", "/api/users",
		))

		allLogs := observedLogs.All()
		if len(allLogs) == 0 {
			t.Fatal("Expected at least one log entry")
		}

		logEntry := allLogs[0]

		// Verify level
		if logEntry.Level != zapcore.ErrorLevel {
			t.Errorf("Expected ERROR level, got %v", logEntry.Level)
		}

		// Verify log_type is error
		foundLogType := false
		for _, field := range logEntry.Context {
			if field.Key == "log_type" && field.String == "error" {
				foundLogType = true
				break
			}
		}

		if !foundLogType {
			t.Error("Expected log_type=error")
		}
	})

	t.Run("should log all log types correctly", func(t *testing.T) {
		tests := []struct {
			name         string
			logFn        func()
			expectedLvl  zapcore.Level
			expectedType string
		}{
			{
				name: "HTTP",
				logFn: func() {
					observedLogs.TakeAll()
					initialized.HTTP(context.Background(), "Request", Fields("method", "GET"))
				},
				expectedLvl:  zapcore.InfoLevel,
				expectedType: "http",
			},
			{
				name: "Security",
				logFn: func() {
					observedLogs.TakeAll()
					initialized.Security(context.Background(), "Security event", Fields("ip", "1.2.3.4"))
				},
				expectedLvl:  zapcore.WarnLevel,
				expectedType: "security",
			},
			{
				name: "Audit",
				logFn: func() {
					observedLogs.TakeAll()
					initialized.Audit(context.Background(), "User action", Fields("user", "admin"))
				},
				expectedLvl:  zapcore.InfoLevel,
				expectedType: "audit",
			},
			{
				name: "Debug",
				logFn: func() {
					observedLogs.TakeAll()
					initialized.Debug(context.Background(), "Debug info", Fields("step", "1"))
				},
				expectedLvl:  zapcore.DebugLevel,
				expectedType: "debug",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.logFn()

				logs := observedLogs.All()
				if len(logs) == 0 {
					t.Fatalf("%s: Expected log entry", tt.name)
				}

				entry := logs[0]

				if entry.Level != tt.expectedLvl {
					t.Errorf("%s: Expected level %v, got %v", tt.name, tt.expectedLvl, entry.Level)
				}

				foundType := false
				for _, field := range entry.Context {
					if field.Key == "log_type" && field.String == tt.expectedType {
						foundType = true
						break
					}
				}

				if !foundType {
					t.Errorf("%s: Expected log_type=%s", tt.name, tt.expectedType)
				}
			})
		}
	})
}

// Example: Alternative approach - return LogEntry for inspection
// This requires modifying the logger to optionally return logged data

type LogEntry struct {
	Level   string
	Message string
	Fields  LogContext
	LogType string
}

// TestableLogger wraps Logger with test mode
type TestableLogger struct {
	*Logger
	testMode   bool
	lastEntry  *LogEntry
	allEntries []LogEntry
}

func (tl *TestableLogger) InfoTestable(message string, logContext LogContext) *LogEntry {
	// Still log normally
	tl.Info(context.Background(), message, logContext)

	// If in test mode, also capture the entry
	if tl.testMode {
		entry := &LogEntry{
			Level:   "INFO",
			Message: message,
			Fields:  make(LogContext),
			LogType: string(TypeNormal),
		}
		// Copy context
		for k, v := range logContext {
			entry.Fields[k] = v
		}
		tl.lastEntry = entry
		tl.allEntries = append(tl.allEntries, *entry)
		return entry
	}

	return nil
}

func TestTestableLogger(t *testing.T) {
	instance = nil
	once = sync.Once{}

	Initialize(Config{
		ServiceName:    "testable-logger-test",
		ServiceVersion: "1.0.0",
		Env:            "test",
		Level:          LevelINFO,
	})

	tl := &TestableLogger{
		Logger:   GetInstance(),
		testMode: true,
	}

	t.Run("should return log entry for inspection", func(t *testing.T) {
		entry := tl.InfoTestable("User logged in", Fields(
			"user_id", "usr-789",
			"ip", "192.168.1.1",
		))

		if entry == nil {
			t.Fatal("Expected log entry to be returned")
		}

		if entry.Message != "User logged in" {
			t.Errorf("Expected message 'User logged in', got '%s'", entry.Message)
		}

		if entry.Level != "INFO" {
			t.Errorf("Expected level INFO, got %s", entry.Level)
		}

		if entry.Fields["user_id"] != "usr-789" {
			t.Errorf("Expected user_id=usr-789, got %v", entry.Fields["user_id"])
		}

		if entry.Fields["ip"] != "192.168.1.1" {
			t.Errorf("Expected ip=192.168.1.1, got %v", entry.Fields["ip"])
		}
	})

	t.Run("should capture all entries", func(t *testing.T) {
		tl.allEntries = []LogEntry{} // Reset

		tl.InfoTestable("First", Fields("order", 1))
		tl.InfoTestable("Second", Fields("order", 2))
		tl.InfoTestable("Third", Fields("order", 3))

		if len(tl.allEntries) != 3 {
			t.Errorf("Expected 3 entries, got %d", len(tl.allEntries))
		}

		for i, entry := range tl.allEntries {
			expectedOrder := i + 1
			if entry.Fields["order"] != expectedOrder {
				t.Errorf("Entry %d: expected order=%d, got %v", i, expectedOrder, entry.Fields["order"])
			}
		}
	})
}

// =============================================================================
// COMPREHENSIVE COVERAGE TESTS (to achieve 100% coverage)
// =============================================================================

func TestGetInstancePanic(t *testing.T) {
	// Reset instance to nil
	instance = nil
	once = sync.Once{}

	// Attempting to get instance without initialization should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected GetInstance to panic when not initialized")
		}
	}()

	GetInstance()
}

func TestAllLogLevels(t *testing.T) {
	instance = nil
	once = sync.Once{}

	tests := []struct {
		name     string
		level    LogLevel
		expected zapcore.Level
	}{
		{"DEBUG level", LevelDEBUG, zapcore.DebugLevel},
		{"INFO level", LevelINFO, zapcore.InfoLevel},
		{"WARN level", LevelWARN, zapcore.WarnLevel},
		{"ERROR level", LevelERROR, zapcore.ErrorLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance = nil
			once = sync.Once{}

			observedCore, _ := observer.New(zapcore.DebugLevel)
			observedLogger := zap.New(observedCore)

			logger := Initialize(Config{
				ServiceName:    "level-test",
				ServiceVersion: "1.0.0",
				Env:            "test",
				Level:          tt.level,
			})

			logger.zap = observedLogger

			// Verify the level conversion
			actualLevel := logger.getZapLevel()
			if actualLevel != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, actualLevel)
			}
		})
	}
}

func TestErrorWithErrorObject(t *testing.T) {
	instance = nil
	once = sync.Once{}

	observedCore, observedLogs := observer.New(zapcore.DebugLevel)
	observedLogger := zap.New(observedCore)

	logger := Initialize(Config{
		ServiceName:    "error-test",
		ServiceVersion: "1.0.0",
		Env:            "test",
		Level:          LevelDEBUG,
	})

	logger.zap = observedLogger

	// Test with actual error object
	testErr := &testError{msg: "test error"}
	logger.Error(context.Background(), "An error occurred", Fields(
		"error", testErr,
		"code", "ERR_500",
	))

	logs := observedLogs.All()
	if len(logs) == 0 {
		t.Fatal("Expected log entry")
	}

	entry := logs[0]

	// Verify error was converted
	hasErrorMessage := false
	hasErrorType := false

	for _, field := range entry.Context {
		if field.Key == "error_message" {
			hasErrorMessage = true
			if field.String != "test error" {
				t.Errorf("Expected error_message='test error', got '%s'", field.String)
			}
		}
		if field.Key == "error_type" {
			hasErrorType = true
		}
	}

	if !hasErrorMessage {
		t.Error("Expected error_message field")
	}
	if !hasErrorType {
		t.Error("Expected error_type field")
	}
}

func TestTraceContext(t *testing.T) {
	instance = nil
	once = sync.Once{}

	observedCore, observedLogs := observer.New(zapcore.DebugLevel)
	observedLogger := zap.New(observedCore)

	logger := Initialize(Config{
		ServiceName:    "trace-test",
		ServiceVersion: "1.0.0",
		Env:            "test",
		Level:          LevelDEBUG,
	})

	logger.zap = observedLogger

	t.Run("should handle nil context", func(t *testing.T) {
		observedLogs.TakeAll()
		logger.Info(nil, "Test with nil context", Fields("test", "value"))

		logs := observedLogs.All()
		if len(logs) == 0 {
			t.Fatal("Expected log entry")
		}

		// Should not have trace_id or span_id fields
		for _, field := range logs[0].Context {
			if field.Key == "trace_id" || field.Key == "span_id" {
				t.Errorf("Should not have trace fields with nil context, found %s", field.Key)
			}
		}
	})

	t.Run("should handle context without trace", func(t *testing.T) {
		observedLogs.TakeAll()
		ctx := context.Background()
		logger.Info(ctx, "Test without trace", Fields("test", "value"))

		logs := observedLogs.All()
		if len(logs) == 0 {
			t.Fatal("Expected log entry")
		}

		// Should not have trace_id or span_id fields
		for _, field := range logs[0].Context {
			if field.Key == "trace_id" || field.Key == "span_id" {
				t.Errorf("Should not have trace fields without span, found %s", field.Key)
			}
		}
	})
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestSync(t *testing.T) {
	instance = nil
	once = sync.Once{}

	logger := Initialize(Config{
		ServiceName:    "sync-test",
		ServiceVersion: "1.0.0",
		Env:            "test",
		Level:          LevelINFO,
	})

	// Sync should not error
	err := logger.Sync()
	if err != nil {
		// Sync can fail on some systems (e.g., /dev/stderr on CI), so we just verify it was called
		t.Logf("Sync returned error (this is acceptable on some systems): %v", err)
	}
}

func TestFieldsValidation(t *testing.T) {
	t.Run("should panic with odd number of arguments", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected Fields to panic with odd number of arguments")
			}
		}()

		Fields("key1", "value1", "key2")
	})

	t.Run("should panic with non-string key", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected Fields to panic with non-string key")
			}
		}()

		Fields(123, "value")
	})

	t.Run("should work with valid arguments", func(t *testing.T) {
		result := Fields(
			"key1", "value1",
			"key2", 123,
			"key3", true,
		)

		if len(result) != 3 {
			t.Errorf("Expected 3 fields, got %d", len(result))
		}

		if result["key1"] != "value1" {
			t.Errorf("Expected key1='value1', got '%v'", result["key1"])
		}

		if result["key2"] != 123 {
			t.Errorf("Expected key2=123, got %v", result["key2"])
		}

		if result["key3"] != true {
			t.Errorf("Expected key3=true, got %v", result["key3"])
		}
	})
}

func TestMeasureDuration(t *testing.T) {
	t.Run("should measure duration correctly", func(t *testing.T) {
		start := time.Now()
		time.Sleep(10 * time.Millisecond)
		duration := MeasureDuration(start)

		if duration < 10 {
			t.Errorf("Expected duration >= 10ms, got %.2fms", duration)
		}

		if duration > 100 {
			t.Errorf("Expected duration < 100ms, got %.2fms", duration)
		}
	})

	t.Run("should return near zero for immediate measurement", func(t *testing.T) {
		start := time.Now()
		duration := MeasureDuration(start)

		if duration > 5 {
			t.Errorf("Expected duration < 5ms, got %.2fms", duration)
		}
	})
}

// Summary: Best practices for testing loggers
//
// Option 1: Observable Logs (RECOMMENDED)
// - Use zap's observer.New() to capture logs
// - Inspect log entries directly without parsing JSON
// - No performance impact on production code
// - Clean separation of concerns
//
// Option 2: Testable Wrapper
// - Wrap logger with test-specific functionality
// - Return LogEntry for inspection
// - Requires wrapper maintenance
// - Good for complex test scenarios
//
// Option 3: Current Approach (stdout capture)
// - Simple but fragile
// - Requires JSON parsing
// - Can be flaky with timing
// - Not recommended for unit tests
//
// Option 4: Integration Tests Only
// - Just verify functions don't panic
// - Fast and simple
// - Good for CI/CD
// - Current approach in main test files
