# logger-go

[![Go Version](https://img.shields.io/badge/go-1.24%2B-blue)](https://golang.org/dl/)
[![GoDoc](https://pkg.go.dev/badge/github.com/rcommerz/logger-go)](https://pkg.go.dev/github.com/rcommerz/logger-go)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rcommerz/logger-go)](https://goreportcard.com/report/github.com/rcommerz/logger-go)
[![Coverage](https://img.shields.io/badge/coverage-99%25-brightgreen)](https://github.com/rcommerz/logger-go)

Production-ready structured logging package for Go microservices with OpenTelemetry support and LGTM stack compatibility.

## Features

- ✅ **Singleton Pattern** - Thread-safe initialization
- ✅ **Structured JSON Logging** - Always outputs valid JSON
- ✅ **OpenTelemetry Integration** - Automatic trace_id and span_id extraction from context
- ✅ **Fiber Middleware** - Automatic HTTP request/response logging
- ✅ **Performance** - Built on Zap (fastest Go logger)
- ✅ **Zero Allocations** - Optimized for high throughput
- ✅ **Container-Friendly** - Outputs to stdout
- ✅ **Production-Ready** - Minimal dependencies
- ✅ **99% Test Coverage** - Thoroughly tested and reliable

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Basic Logging](#basic-logging)
  - [Fiber Middleware](#fiber-middleware)
  - [OpenTelemetry Integration](#opentelemetry-integration)
- [API Reference](#api-reference)
- [Testing](#testing)
- [Best Practices](#best-practices)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/rcommerz/logger-go
```

**Requirements:**
- Go 1.24 or higher

## Quick Start

### 1. Initialize Logger (Do this once in main.go)

```go
package main

import (
    "github.com/rcommerz/logger-go"
)

func main() {
    // Initialize logger once
    logger.Initialize(logger.Config{
        ServiceName:    "product-service",
        ServiceVersion: "1.2.0",
        Env:            "production",
        Level:          "info", // debug, info, warn, error
    })
    
    // Always sync on shutdown
    defer logger.GetInstance().Sync()
    
    // Your application code...
}
```

### 2. Use Logger Anywhere

```go
package services

import (
    "context"
    "github.com/rcommerz/logger-go"
)

func ProcessOrder(ctx context.Context, orderID string) {
    log := logger.GetInstance()
    
    // Basic logging (pass context for OpenTelemetry trace extraction)
    log.Info(ctx, "Processing order", logger.Fields(
        "order_id", orderID,
        "status", "pending",
    ))
    
    // Error logging
    if err := validateOrder(orderID); err != nil {
        log.Error(ctx, "Order validation failed", logger.Fields(
            "order_id", orderID,
            "error", err.Error(),
        ))
        return
    }
    
    // Warning
    log.Warn(ctx, "Low inventory detected", logger.Fields(
        "product_id", "PROD-123",
        "quantity", 5,
    ))
    
    // Debug (only logs if level is debug)
    log.Debug(ctx, "Order processing steps", logger.Fields(
        "step", "payment_validation",
        "order_id", orderID,
    ))
    
    // Security logging
    log.Security(ctx, "Unauthorized access attempt", logger.Fields(
        "user_id", "USR-456",
        "resource", "/admin/users",
        "ip", "192.168.1.100",
    ))
    
    // Audit logging
    log.Audit(ctx, "User role changed", logger.Fields(
        "user_id", "USR-789",
        "admin_id", "ADM-001",
        "old_role", "user",
        "new_role", "admin",
    ))
}
```

### 3. Fiber HTTP Middleware

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/rcommerz/logger-go"
)

func main() {
    // Initialize logger
    logger.Initialize(logger.Config{
        ServiceName:    "api-gateway",
        ServiceVersion: "1.0.0",
        Env:            "production",
        Level:          "info",
    })
    defer logger.GetInstance().Sync()
    
    // Create Fiber app
    app := fiber.New()
    
    // Add logger middleware (automatically extracts OpenTelemetry trace from request context)
    app.Use(logger.FiberMiddleware(&logger.MiddlewareOptions{
        ExcludePaths:   []string{"/health", "/metrics"},
        IncludeHeaders: false,
    }))
    
    // Add recovery middleware (logs panics)
    app.Use(logger.RecoveryMiddleware())
    
    // Your routes
    app.Get("/api/products", getProducts)
    app.Post("/api/orders", createOrder)
    
    app.Listen(":8080")
}

func getProducts(c *fiber.Ctx) error {
    log := logger.GetInstance()
    
    // Use c.UserContext() to pass the request context (includes OpenTelemetry trace)
    log.Info(c.UserContext(), "Fetching products", logger.Fields(
        "user_id", c.Locals("user_id"),
        "category", c.Query("category"),
    ))
    
    return c.JSON(fiber.Map{"products": []string{}})
}
```

### 4. OpenTelemetry Integration

```go
package main

import (
    "context"
    
    "github.com/rcommerz/logger-go"
    "go.opentelemetry.io/otel"
)

func processPayment(ctx context.Context, paymentID string) error {
    tracer := otel.Tracer("payment-service")
    ctx, span := tracer.Start(ctx, "process-payment")
    defer span.End()
    
    log := logger.GetInstance()
    
    // Logger automatically extracts trace_id and span_id from context
    log.Info(ctx, "Processing payment", logger.Fields(
        "payment_id", paymentID,
        "amount", 199.99,
    ))
    
    // If payment fails
    err := chargeCard(paymentID)
    if err != nil {
        log.Error(ctx, "Payment failed", logger.Fields(
            "payment_id", paymentID,
            "error", err.Error(),
        ))
        span.RecordError(err)
        return err
    }
    
    return nil
}
```

## Example Output

### Standard Log

```json
{
  "@timestamp": "2026-02-22T10:30:45.123Z",
  "log.level": "INFO",
  "log_type": "normal",
  "service.name": "product-service",
  "service.version": "1.2.0",
  "env": "production",
  "host.name": "pod-product-abc123",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "message": "Order created",
  "order_id": "ORD-12345",
  "user_id": "USR-999",
  "total_amount": 199.99,
  "items_count": 3
}
```

### HTTP Request Log

```json
{
  "@timestamp": "2026-02-22T10:30:46.456Z",
  "log.level": "INFO",
  "log_type": "http",
  "service.name": "api-gateway",
  "service.version": "1.0.0",
  "env": "production",
  "host.name": "pod-gateway-xyz789",
  "trace_id": "3ad45f8b21c34e2a8d41ba6e3f9c0412",
  "span_id": "91f23ab4cd5e6789",
  "message": "GET /api/products 200",
  "method": "GET",
  "path": "/api/products",
  "status_code": 200,
  "duration_ms": 45.2,
  "client_ip": "10.0.1.25",
  "user_agent": "Go-http-client/1.1"
}
```

### Error Log

```json
{
  "@timestamp": "2026-02-22T10:30:47.789Z",
  "log.level": "ERROR",
  "log_type": "error",
  "service.name": "payment-service",
  "service.version": "2.1.0",
  "env": "production",
  "host.name": "pod-payment-def456",
  "trace_id": "7cd89e12f43b4a9f8e21dc7a5b6f3028",
  "span_id": "12a34b56c78d90ef",
  "message": "Payment processing failed",
  "payment_id": "PAY-67890",
  "error": "insufficient funds",
  "user_id": "USR-111",
  "amount": 500.00
}
```

## API Reference

### Logger Methods

All logging methods now require a `context.Context` as the first parameter for OpenTelemetry trace extraction.

#### `Info(ctx context.Context, message string, fields LogContext)`

Log informational messages.

#### `Error(ctx context.Context, message string, fields LogContext)`

Log error messages. Automatically extracts error type from error objects.

#### `Warn(ctx context.Context, message string, fields LogContext)`

Log warning messages.

#### `Debug(ctx context.Context, message string, fields LogContext)`

Log debug messages (only if level is debug).

#### `Security(ctx context.Context, message string, fields LogContext)`

Log security-related events (log_type = "security").

#### `Audit(ctx context.Context, message string, fields LogContext)`

Log audit trail events (log_type = "audit").

#### `HTTP(ctx context.Context, message string, fields LogContext)`

Log HTTP-specific events (log_type = "http").

### Helper Functions

#### `Fields(keyValues ...interface{}) LogContext`

Helper to create fields map:

```go
logger.Fields("key1", "value1", "key2", 123)
```

#### `MeasureDuration(start time.Time) float64`

Calculate duration in milliseconds:

```go
start := time.Now()
// ... do work ...
duration := logger.MeasureDuration(start)
```

### Middleware Options

#### `FiberMiddleware(options *MiddlewareOptions) fiber.Handler`

- `ExcludePaths []string` - Paths to exclude from logging
- `IncludeHeaders bool` - Include request headers (default: false)

**Note**: The middleware automatically passes `c.UserContext()` to the logger, enabling automatic OpenTelemetry trace extraction.

#### `RecoveryMiddleware() fiber.Handler`

Middleware that recovers from panics and logs them with full context and trace information.

## Best Practices

### ✅ DO

- Initialize logger once in `main()`
- Always call `defer logger.GetInstance().Sync()` on shutdown
- Use `Fields()` helper for readability
- Pass OpenTelemetry spans when available
- Use appropriate log levels
- Set log level via environment variable

### ❌ DON'T

- Initialize logger multiple times
- Log sensitive information (passwords, tokens, credit cards)
- Use debug level in production
- Forget to call `Sync()` on shutdown

## Environment Variables

```bash
# Set via Config.Level
export LOG_LEVEL=info  # debug, info, warn, error
```

## LGTM Stack Compatibility

Designed for:

- **Loki** - Structured JSON with consistent labels
- **Grafana** - Standardized field names for dashboards
- **Tempo** - Automatic trace correlation
- **Prometheus** - Derive metrics from logs

### Example Loki Query

```logql
{service_name="product-service"} | json | log_type="http" | status_code >= 500
```

## Performance

Built on Zap, one of the fastest structured loggers for Go:

- Zero allocations in hot paths
- Optimized for throughput
- Minimal CPU overhead

## Testing Your Application

### Why Logger Doesn't Return Data

The logger uses `zap.Logger` which writes directly to **stdout** for maximum performance. The methods (`Info`, `Error`, etc.) are **void functions** by design—they don't return values because:

1. **Performance**: No overhead of creating return objects
2. **Design**: Follows Zap's zero-allocation pattern
3. **Production**: Logs are written once, not captured

### Testing Strategies

#### ✅ **Recommended: Observable Logs**

Use Zap's built-in observer for **direct log inspection** (no stdout parsing):

```go
import (
    "testing"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "go.uber.org/zap/zaptest/observer"
    "github.com/rcommerz/logger-go"
)

func TestMyService(t *testing.T) {
    // Create observable core
    observedCore, observedLogs := observer.New(zapcore.DebugLevel)
    observedLogger := zap.New(observedCore)
    
    // Initialize and replace zap instance
    log := logger.Initialize(logger.Config{
        ServiceName:    "test-service",
        ServiceVersion: "1.0.0",
        Env:            "test",
        Level:          logger.LevelDEBUG,
    })
    log.zap = observedLogger
    
    // Run your code
    log.Info(context.Background(), "User action", logger.Fields("user_id", "123"))
    
    // Inspect logs directly (NO JSON parsing!)
    allLogs := observedLogs.All()
    entry := allLogs[0]
    
    assert.Equal(t, "User action", entry.Message)
    assert.Equal(t, zapcore.InfoLevel, entry.Level)
    
    // Check fields
    for _, field := range entry.Context {
        if field.Key == "user_id" {
            assert.Equal(t, "123", field.String)
        }
    }
}
```

**Benefits:**

- ✅ Direct access to log entries
- ✅ No stdout capture needed
- ✅ No JSON parsing required
- ✅ Fast and reliable
- ✅ Type-safe field inspection

#### ✅ **Alternative: Integration Tests**

For simpler tests, just verify no panics occur:

```go
func TestServiceLogging(t *testing.T) {
    logger.Initialize(logger.Config{
        ServiceName: "test",
        ServiceVersion: "1.0.0",
        Env: "test",
        Level: logger.LevelINFO,
    })
    
    log := logger.GetInstance()
    ctx := context.Background()
    
    // Just verify these don't panic
    log.Info(ctx, "Test message", logger.Fields("key", "value"))
    log.Error(ctx, "Error message", logger.Fields("error", "test"))
    
    // Test passes if no panic
}
```

**Benefits:**

- ✅ Very fast
- ✅ Simple to write
- ✅ Good for CI/CD
- ✅ Tests the actual behavior

#### ⚠️ **Not Recommended: Stdout Capture**

Avoid capturing and parsing stdout—it's fragile and slow:

```go
// ❌ Don't do this - flaky and slow
func captureStdout() string {
    // Captures os.Stdout, parses JSON, brittle...
}
```

### Test Files in This Package

- **`logger_test.go`** - Consolidated logger tests (100% coverage)
- **`middleware_test.go`** - Consolidated middleware tests (100% coverage)
- **`logger_testable_test.go`** - Examples of observable and wrapper patterns

### Running Tests

```bash
# Run all tests
go test -v

# Run with coverage (100% achieved!)
go test -cover

# Run specific test
go test -run TestLoggingWithObserver

# Generate HTML coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Coverage

Current coverage: **99%** of statements ✅

All core functionality is tested:

- Logger initialization and singleton pattern
- All log levels (DEBUG, INFO, WARN, ERROR)
- All log types (normal, http, error, security, audit, debug)
- Error handling and panic recovery
- OpenTelemetry trace context extraction
- Fiber middleware (request/response logging)
- Recovery middleware (panic recovery)
- Fields validation
- Concurrency safety
- Performance characteristics

## Performance

Built on top of [Zap](https://github.com/uber-go/zap), the fastest structured logger for Go:

- **Zero allocations** for most log operations
- **Microsecond latency** for structured logging
- **High throughput** - handles millions of logs per second
- **Minimal CPU overhead** - optimized for production workloads

## Comparison with Other Loggers

| Feature | logger-go | logrus | zap | zerolog |
|---------|-----------|--------|-----|---------|
| Performance | ⚡ Very Fast | Slow | ⚡ Very Fast | ⚡ Very Fast |
| OpenTelemetry | ✅ Built-in | ❌ Manual | ❌ Manual | ❌ Manual |
| Fiber Support | ✅ Native | ⚠️ Custom | ⚠️ Custom | ⚠️ Custom |
| Type Safety | ✅ Yes | ⚠️ Partial | ✅ Yes | ✅ Yes |
| Singleton | ✅ Built-in | ❌ Manual | ❌ Manual | ❌ Manual |
| JSON Output | ✅ Always | ✅ Optional | ✅ Optional | ✅ Always |

## Roadmap

- [ ] Support for additional web frameworks (Gin, Echo)
- [ ] Custom log formatters
- [ ] Log sampling for high-volume scenarios
- [ ] Integration with more observability platforms
- [ ] Performance benchmarks and optimizations

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests (maintain 95%+ coverage)
5. Run tests (`go test -v`)
6. Commit your changes (`git commit -m 'feat: add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and release notes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support & Contact

- **Issues**: [GitHub Issues](https://github.com/rcommerz/logger-go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/rcommerz/logger-go/discussions)
- **Security**: Report vulnerabilities via GitHub Security Advisories

## Acknowledgments

- Built on top of [Uber's Zap](https://github.com/uber-go/zap)
- Inspired by observability best practices
- Thanks to all [contributors](https://github.com/rcommerz/logger-go/graphs/contributors)

## Related Projects

- [logger-laravel](https://github.com/rcommerz/logger-laravel) - PHP/Laravel logging package
- [logger-express](https://github.com/rcommerz/logger-express) - Node.js/Express logging package

---

Made with ❤️ by the RCOMMERZ team
