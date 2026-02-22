package logger

import (
	"context"
	"os"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *Logger
	once     sync.Once
)

// Logger is a structured logger wrapper around zap
type Logger struct {
	zap    *zap.Logger
	config Config
}

// Initialize creates and returns a singleton logger instance
func Initialize(config Config) *Logger {
	once.Do(func() {
		instance = &Logger{
			config: config,
		}
		instance.zap = instance.buildZapLogger()
	})
	return instance
}

// GetInstance returns the singleton logger instance
func GetInstance() *Logger {
	if instance == nil {
		panic("Logger not initialized. Call Initialize() first.")
	}
	return instance
}

// buildZapLogger creates a configured zap logger
func (l *Logger) buildZapLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "log.level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	hostname, _ := os.Hostname()

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		l.getZapLevel(),
	)

	logger := zap.New(core)

	// Add constant fields
	logger = logger.With(
		zap.String("service.name", l.config.ServiceName),
		zap.String("service.version", l.config.ServiceVersion),
		zap.String("env", l.config.Env),
		zap.String("host.name", hostname),
	)

	return logger
}

// getZapLevel converts LogLevel to zapcore.Level
func (l *Logger) getZapLevel() zapcore.Level {
	switch l.config.Level {
	case LevelDEBUG:
		return zapcore.DebugLevel
	case LevelWARN:
		return zapcore.WarnLevel
	case LevelERROR:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// getTraceContext extracts trace_id and span_id from OpenTelemetry context
func (l *Logger) getTraceContext(ctx context.Context) []zap.Field {
	if ctx == nil {
		return []zap.Field{}
	}

	spanContext := trace.SpanFromContext(ctx).SpanContext()
	if !spanContext.IsValid() {
		return []zap.Field{}
	}

	return []zap.Field{
		zap.String("trace_id", spanContext.TraceID().String()),
		zap.String("span_id", spanContext.SpanID().String()),
	}
}

// buildFields converts LogContext to zap.Field array
func (l *Logger) buildFields(ctx context.Context, logType LogType, context LogContext) []zap.Field {
	fields := []zap.Field{
		zap.String("log_type", string(logType)),
	}

	// Add trace context
	fields = append(fields, l.getTraceContext(ctx)...)

	// Add custom context fields
	for key, value := range context {
		fields = append(fields, zap.Any(key, value))
	}

	return fields
}

// Info logs an informational message
func (l *Logger) Info(ctx context.Context, message string, context LogContext) {
	fields := l.buildFields(ctx, TypeNormal, context)
	l.zap.Info(message, fields...)
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, message string, context LogContext) {
	// Handle error objects
	if err, ok := context["error"].(error); ok {
		context["error_message"] = err.Error()
		context["error_type"] = "error"
		delete(context, "error")
	}

	fields := l.buildFields(ctx, TypeError, context)
	l.zap.Error(message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, message string, context LogContext) {
	fields := l.buildFields(ctx, TypeNormal, context)
	l.zap.Warn(message, fields...)
}

// Debug logs a debug message
func (l *Logger) Debug(ctx context.Context, message string, context LogContext) {
	fields := l.buildFields(ctx, TypeDebug, context)
	l.zap.Debug(message, fields...)
}

// HTTP logs an HTTP request/response
func (l *Logger) HTTP(ctx context.Context, message string, context LogContext) {
	fields := l.buildFields(ctx, TypeHTTP, context)
	l.zap.Info(message, fields...)
}

// Security logs a security-related event
func (l *Logger) Security(ctx context.Context, message string, context LogContext) {
	fields := l.buildFields(ctx, TypeSecurity, context)
	l.zap.Warn(message, fields...)
}

// Audit logs an audit trail event
func (l *Logger) Audit(ctx context.Context, message string, context LogContext) {
	fields := l.buildFields(ctx, TypeAudit, context)
	l.zap.Info(message, fields...)
}

// Sync flushes any buffered log entries (call before app shutdown)
func (l *Logger) Sync() error {
	return l.zap.Sync()
}
