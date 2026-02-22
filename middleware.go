package logger

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MiddlewareOptions configures the HTTP logging middleware
type MiddlewareOptions struct {
	ExcludePaths   []string
	IncludeHeaders bool
	IncludeBody    bool
}

// FiberMiddleware returns a Fiber middleware that logs HTTP requests
func FiberMiddleware(opts *MiddlewareOptions) fiber.Handler {
	if opts == nil {
		opts = &MiddlewareOptions{}
	}

	logger := GetInstance()

	return func(c *fiber.Ctx) error {
		// Skip excluded paths
		path := c.Path()
		for _, excludePath := range opts.ExcludePaths {
			if path == excludePath {
				return c.Next()
			}
		}

		startTime := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Build log context
		context := LogContext{
			"method":      c.Method(),
			"path":        path,
			"status_code": c.Response().StatusCode(),
			"duration_ms": duration.Milliseconds(),
			"ip":          c.IP(),
			"user_agent":  c.Get("User-Agent"),
		}

		// Add query params if present
		if len(c.Context().QueryArgs().String()) > 0 {
			context["query"] = c.Context().QueryArgs().String()
		}

		// Add headers if requested
		if opts.IncludeHeaders {
			headers := make(map[string]string)
			c.Request().Header.VisitAll(func(key, value []byte) {
				headers[string(key)] = string(value)
			})
			context["headers"] = headers
		}

		// Add user_id from locals if available
		if userID := c.Locals("user_id"); userID != nil {
			context["user_id"] = userID
		}

		// Build message
		message := fmt.Sprintf("%s %s %d", c.Method(), path, c.Response().StatusCode())

		// Log based on status code
		statusCode := c.Response().StatusCode()
		ctx := c.UserContext()
		if statusCode >= 500 {
			logger.Error(ctx, message, context)
		} else if statusCode >= 400 {
			logger.Warn(ctx, message, context)
		} else {
			logger.HTTP(ctx, message, context)
		}

		return err
	}
}

// RecoveryMiddleware returns a Fiber middleware that recovers from panics and logs them
func RecoveryMiddleware() fiber.Handler {
	logger := GetInstance()

	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				context := LogContext{
					"method":      c.Method(),
					"path":        c.Path(),
					"panic":       r,
					"status_code": 500,
				}

				logger.Error(c.UserContext(), "Panic recovered", context)
				err = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "internal server error",
				})
			}
		}()

		return c.Next()
	}
}
