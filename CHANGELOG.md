# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-02-23

### Added

- Initial release of logger-go
- Structured JSON logging with Zap
- OpenTelemetry integration with automatic trace_id and span_id extraction
- Fiber middleware for HTTP request/response logging
- Recovery middleware for panic handling
- Support for multiple log types: Info, Error, Warn, Debug, HTTP, Security, Audit
- Thread-safe singleton pattern
- Context-aware logging for distributed tracing
- Helper functions: Fields() and MeasureDuration()
- Comprehensive test suite with 99% coverage
- Production-ready configuration

### Features

- Singleton pattern with thread-safe initialization
- Structured JSON logging
- OpenTelemetry trace context extraction
- Fiber HTTP middleware
- Zero-allocation optimizations
- Container-friendly stdout logging
- Minimal dependencies
- LGTM stack compatible

[1.0.0]: https://github.com/rcommerz/logger-go/releases/tag/v1.0.0
