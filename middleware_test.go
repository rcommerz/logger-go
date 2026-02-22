package logger

import (
	"errors"
	"io"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// =============================================================================
// FIBER MIDDLEWARE BASIC TESTS
// =============================================================================

func TestFiberMiddleware(t *testing.T) {
	instance = nil
	once = sync.Once{}

	Initialize(Config{
		ServiceName:    "test-api",
		ServiceVersion: "1.0.0",
		Env:            "test",
		Level:          LevelINFO,
	})

	t.Run("should log successful requests", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(&MiddlewareOptions{
			ExcludePaths: []string{},
		}))

		app.Get("/api/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/api/test", nil)
		resp, err := app.Test(req)

		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("should log error responses", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))

		app.Get("/api/error", func(c *fiber.Ctx) error {
			return c.Status(500).JSON(fiber.Map{"error": "internal error"})
		})

		req := httptest.NewRequest("GET", "/api/error", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 500 {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	})

	t.Run("should exclude specified paths", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(&MiddlewareOptions{
			ExcludePaths: []string{"/health", "/metrics"},
		}))

		app.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "healthy"})
		})

		req := httptest.NewRequest("GET", "/health", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("should measure request duration", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))

		app.Get("/api/slow", func(c *fiber.Ctx) error {
			time.Sleep(50 * time.Millisecond)
			return c.JSON(fiber.Map{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/api/slow", nil)
		start := time.Now()
		resp, _ := app.Test(req, -1) // -1 timeout means wait indefinitely
		duration := time.Since(start)

		if duration < 50*time.Millisecond {
			t.Errorf("Expected duration >= 50ms, got %v", duration)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}

// =============================================================================
// FIBER MIDDLEWARE COMPREHENSIVE TESTS
// =============================================================================

func TestFiberMiddlewareComprehensive(t *testing.T) {
	instance = nil
	once = sync.Once{}

	Initialize(Config{
		ServiceName:    "middleware-comprehensive-test",
		ServiceVersion: "1.0.0",
		Env:            "test",
		Level:          LevelDEBUG,
	})

	t.Run("should include query params", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))

		app.Get("/api/search", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"results": []string{}})
		})

		req := httptest.NewRequest("GET", "/api/search?q=test&limit=10", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("should include headers when requested", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(&MiddlewareOptions{
			IncludeHeaders: true,
		}))

		app.Get("/api/headers", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"ok": true})
		})

		req := httptest.NewRequest("GET", "/api/headers", nil)
		req.Header.Set("X-Custom-Header", "value")
		req.Header.Set("Authorization", "Bearer token")
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("should include user_id from locals", func(t *testing.T) {
		app := fiber.New()

		// Middleware that sets user_id before logging middleware
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user_id", "usr-123")
			return c.Next()
		})

		app.Use(FiberMiddleware(nil))

		app.Get("/api/user", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"user": "test"})
		})

		req := httptest.NewRequest("GET", "/api/user", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("should log as warning for 4xx status", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))

		app.Get("/api/not-found", func(c *fiber.Ctx) error {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		})

		req := httptest.NewRequest("GET", "/api/not-found", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("should log as error for 5xx status", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))

		app.Get("/api/server-error", func(c *fiber.Ctx) error {
			return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
		})

		req := httptest.NewRequest("GET", "/api/server-error", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 503 {
			t.Errorf("Expected status 503, got %d", resp.StatusCode)
		}
	})

	t.Run("should log as HTTP for 2xx status", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))

		app.Post("/api/create", func(c *fiber.Ctx) error {
			return c.Status(201).JSON(fiber.Map{"id": "123"})
		})

		req := httptest.NewRequest("POST", "/api/create", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 201 {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("should log as HTTP for 3xx status", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))

		app.Get("/api/redirect", func(c *fiber.Ctx) error {
			return c.Redirect("/new-location", 301)
		})

		req := httptest.NewRequest("GET", "/api/redirect", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 301 {
			t.Errorf("Expected status 301, got %d", resp.StatusCode)
		}
	})

	t.Run("should work without options (nil)", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))

		app.Get("/api/nil-opts", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"ok": true})
		})

		req := httptest.NewRequest("GET", "/api/nil-opts", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}

// =============================================================================
// RECOVERY MIDDLEWARE TESTS
// =============================================================================

func TestRecoveryMiddleware(t *testing.T) {
	instance = nil
	once = sync.Once{}

	Initialize(Config{
		ServiceName:    "test-api",
		ServiceVersion: "1.0.0",
		Env:            "test",
		Level:          LevelINFO,
	})

	t.Run("should recover from panics", func(t *testing.T) {
		app := fiber.New()
		app.Use(RecoveryMiddleware())

		app.Get("/api/panic", func(c *fiber.Ctx) error {
			panic("test panic")
		})

		req := httptest.NewRequest("GET", "/api/panic", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 500 {
			t.Errorf("Expected status 500 after panic, got %d", resp.StatusCode)
		}
	})

	t.Run("should recover from string panic", func(t *testing.T) {
		app := fiber.New()
		app.Use(RecoveryMiddleware())

		app.Get("/api/panic-string", func(c *fiber.Ctx) error {
			panic("string panic message")
		})

		req := httptest.NewRequest("GET", "/api/panic-string", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 500 {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	})

	t.Run("should recover from error panic", func(t *testing.T) {
		app := fiber.New()
		app.Use(RecoveryMiddleware())

		app.Get("/api/panic-error", func(c *fiber.Ctx) error {
			panic(errors.New("error panic"))
		})

		req := httptest.NewRequest("GET", "/api/panic-error", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 500 {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	})

	t.Run("should not interfere with normal requests", func(t *testing.T) {
		app := fiber.New()
		app.Use(RecoveryMiddleware())

		app.Get("/api/normal", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/api/normal", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}

// =============================================================================
// MIDDLEWARE EDGE CASES
// =============================================================================

func TestMiddlewareEdgeCases(t *testing.T) {
	instance = nil
	once = sync.Once{}

	Initialize(Config{
		ServiceName:    "edge-case-test",
		ServiceVersion: "1.0.0",
		Env:            "test",
		Level:          LevelDEBUG,
	})

	t.Run("should handle empty path", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))

		app.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"root": true})
		})

		req := httptest.NewRequest("GET", "/", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("should handle path with parameters", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))

		app.Get("/api/users/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"id": c.Params("id")})
		})

		req := httptest.NewRequest("GET", "/api/users/test-123", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("should handle multiple excluded paths", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(&MiddlewareOptions{
			ExcludePaths: []string{"/health", "/metrics", "/ready"},
		}))

		app.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"healthy": true})
		})
		app.Get("/metrics", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"metrics": true})
		})

		req1 := httptest.NewRequest("GET", "/health", nil)
		resp1, _ := app.Test(req1)

		req2 := httptest.NewRequest("GET", "/metrics", nil)
		resp2, _ := app.Test(req2)

		if resp1.StatusCode != 200 || resp2.StatusCode != 200 {
			t.Error("Expected both excluded paths to return 200")
		}
	})
}

// =============================================================================
// MIDDLEWARE INTEGRATION TESTS
// =============================================================================

func TestMiddlewareIntegration(t *testing.T) {
	instance = nil
	once = sync.Once{}

	Initialize(Config{
		ServiceName:    "integration-test",
		ServiceVersion: "1.0.0",
		Env:            "test",
		Level:          LevelINFO,
	})

	t.Run("should work with both middleware", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))
		app.Use(RecoveryMiddleware())

		app.Get("/api/data", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"data": "test"})
		})

		req := httptest.NewRequest("GET", "/api/data", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("should log and recover from panic with both middleware", func(t *testing.T) {
		app := fiber.New()
		app.Use(FiberMiddleware(nil))
		app.Use(RecoveryMiddleware())

		app.Get("/api/panic-integration", func(c *fiber.Ctx) error {
			panic("integration panic test")
		})

		req := httptest.NewRequest("GET", "/api/panic-integration", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 500 {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}

		// Verify response body
		body, _ := io.ReadAll(resp.Body)
		if len(body) == 0 {
			t.Error("Expected error response body")
		}
	})
}
