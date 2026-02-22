# Contributing to logger-go

Thank you for your interest in contributing to logger-go! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and constructive in all interactions.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:

- Clear description of the bug
- Steps to reproduce
- Expected vs actual behavior
- Go version and OS
- Relevant code snippets or logs

### Suggesting Features

Feature suggestions are welcome! Please:

- Check if the feature already exists or is planned
- Provide a clear use case
- Explain why it would benefit users

### Pull Requests

1. **Fork the repository**

   ```bash
   git clone https://github.com/rcommerz/logger-go.git
   cd logger-go
   ```

2. **Create a feature branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Write clean, documented code
   - Follow existing code style
   - Add tests for new functionality
   - Ensure all tests pass: `go test -v`

4. **Run tests and checks**

   ```bash
   # Run all tests
   go test -v

   # Check test coverage
   go test -cover

   # Format code
   go fmt ./...

   # Run linter (if available)
   golangci-lint run
   ```

5. **Commit your changes**

   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

   Follow [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` New feature
   - `fix:` Bug fix
   - `docs:` Documentation changes
   - `test:` Test additions or changes
   - `refactor:` Code refactoring
   - `perf:` Performance improvements

6. **Push and create PR**

   ```bash
   git push origin feature/your-feature-name
   ```

   Create a pull request on GitHub with:
   - Clear description of changes
   - Reference to related issues
   - Test results

## Development Guidelines

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions small and focused
- Write self-documenting code
- Add comments for complex logic

### Testing

- Maintain test coverage above 95%
- Write unit tests for all new code
- Use table-driven tests where appropriate
- Test edge cases and error conditions

### Documentation

- Update README.md for user-facing changes
- Add inline comments for complex code
- Update CHANGELOG.md
- Include examples for new features

## Questions?

Feel free to open an issue for any questions or clarifications.

Thank you for contributing! 🎉
