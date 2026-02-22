# Publishing logger-go Package

This guide explains how to publish and maintain the logger-go package using automated GitHub Actions workflows.

## Overview

The logger-go package is published to the Go module ecosystem through GitHub releases. When you create a release with a version tag, the package becomes automatically available via `go get`.

## Prerequisites

1. **GitHub Repository**
   - Repository should be public at `github.com/rcommerz/logger-go`
   - Repository must contain all source files
   - GitHub Actions enabled (Settings → Actions → Allow all actions)

2. **Go Module**
   - `go.mod` file with correct module path: `github.com/rcommerz/logger-go`
   - All dependencies properly declared and tidied

3. **Test Coverage**
   - Minimum 95% test coverage required
   - All tests must pass

4. **Required Files**
   - `README.md` - Package documentation
   - `LICENSE` - MIT license
   - `CHANGELOG.md` - Version history
   - `go.mod` - Go module definition

## Automated Publishing (Recommended)

The package uses GitHub Actions for automated testing and publishing. Two workflows are available:

### Option 1: Automatic Release (Push Tag)

**Best for**: Quick releases when everything is ready

1. **Prepare your changes:**
   ```bash
   # Make your changes
   git add .
   git commit -m "feat: add new feature"
   
   # Update CHANGELOG.md with new version
   # Update version references in documentation
   ```

2. **Create and push a version tag:**
   ```bash
   # Create annotated tag (replace X.Y.Z with your version)
   git tag -a v1.0.0 -m "Release v1.0.0 - Initial stable release"
   
   # Push tag to GitHub
   git push origin v1.0.0
   ```

3. **GitHub Actions automatically:**
   - ✅ Runs all tests across Go 1.21, 1.22, 1.23
   - ✅ Verifies 95%+ test coverage
   - ✅ Validates go.mod is up to date
   - ✅ Checks required files exist
   - ✅ Creates GitHub Release with changelog
   - ✅ Publishes to pkg.go.dev (automatic)

### Option 2: Manual Trigger (Workflow Dispatch)

**Best for**: Creating releases from GitHub UI with version input

1. **Go to GitHub Actions:**
   - Navigate to: `https://github.com/rcommerz/logger-go/actions`
   - Select "Release Package" workflow
   - Click "Run workflow"

2. **Fill in the form:**
   - **Branch**: `master` or `main`
   - **Version**: `1.0.0` (without 'v' prefix)
   - **Create git tag**: ✅ (checked)

3. **Click "Run workflow"**

4. **GitHub Actions automatically:**
   - ✅ Validates version format (X.Y.Z)
   - ✅ Checks if tag already exists
   - ✅ Runs all tests with coverage verification
   - ✅ Creates and pushes git tag
   - ✅ Creates GitHub Release
   - ✅ Publishes to pkg.go.dev

## Version Management

Follow [Semantic Versioning](https://semver.org/):

- **Major** (v2.0.0): Breaking changes
  - Example: Changed function signatures, removed public APIs
  - Increment when: API changes break backward compatibility

- **Minor** (v1.1.0): New features, backward compatible
  - Example: Added new logging methods, new middleware options
  - Increment when: Adding functionality without breaking existing code

- **Patch** (v1.0.1): Bug fixes, backward compatible
  - Example: Fixed race condition, corrected log formatting
  - Increment when: Bug fixes that don't add features

### Version Examples

```bash
# Bug fix release
git tag -a v1.0.1 -m "Fix: race condition in logger initialization"

# New feature release
git tag -a v1.1.0 -m "Feature: add Debug level logging support"

# Breaking change release
git tag -a v2.0.0 -m "Breaking: require context.Context in all logging methods"
```

## Pre-Release Checklist

Before creating a release, ensure:

- [ ] All tests pass: `go test -v ./...`
- [ ] Coverage is ≥95%: `go test -cover`
- [ ] Code is formatted: `go fmt ./...`
- [ ] Dependencies are tidy: `go mod tidy`
- [ ] CHANGELOG.md is updated
- [ ] README.md reflects new changes
- [ ] Version number follows semver
- [ ] Breaking changes are documented

## Updating CHANGELOG.md

Follow [Keep a Changelog](https://keepachangelog.com/) format:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- New feature X

### Changed
- Modified behavior Y

### Fixed
- Bug fix Z

## [1.0.0] - 2026-02-23

### Added
- Initial stable release
- OpenTelemetry integration
- Fiber middleware support
- 99% test coverage
```

## Manual Publishing Steps

If you need to publish manually without GitHub Actions:

### 1. Prepare the Release

```bash
# Run tests
go test -v -cover ./...

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy

# Verify module
go mod verify

# Check coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### 2. Create Git Repository (First Time Only)

```bash
cd /path/to/logger-go

# Initialize git (if not already done)
git init

# Add all files
git add .

# Commit
git commit -m "feat: initial release v1.0.0"

# Add remote
git remote add origin https://github.com/rcommerz/logger-go.git

# Push to GitHub
git push -u origin master
```

### 3. Create Version Tag

```bash
# Create annotated tag
git tag -a v1.0.0 -m "Release v1.0.0 - Initial stable release"

# Push tag to GitHub
git push origin v1.0.0
```

### 4. Create GitHub Release

1. Go to: `https://github.com/rcommerz/logger-go/releases/new`
2. Select tag: `v1.0.0`
3. Release title: `Release v1.0.0`
4. Description: Copy from CHANGELOG.md
5. Check "Set as the latest release"
6. Click "Publish release"

### 5. Verify Publication

The package becomes available immediately:

```bash
# Fetch the package
go get github.com/rcommerz/logger-go@v1.0.0

# View on pkg.go.dev
open https://pkg.go.dev/github.com/rcommerz/logger-go@v1.0.0
```

## Troubleshooting

### Package Not Found

If `go get` fails with "package not found":

1. **Check tag exists:**
   ```bash
   git tag
   git push origin --tags
   ```

2. **Verify repository is public**

3. **Force Go proxy refresh:**
   ```bash
   GOPROXY=https://proxy.golang.org GO111MODULE=on \
     go get github.com/rcommerz/logger-go@v1.0.0
   ```

4. **Request indexing manually:**
   Visit: `https://pkg.go.dev/github.com/rcommerz/logger-go@v1.0.0`

### GitHub Actions Failed

If the release workflow fails:

1. **Check workflow logs:**
   Go to Actions tab → Select failed workflow → View logs

2. **Common issues:**
   - Test coverage below 95%: Add more tests
   - go.mod not tidy: Run `go mod tidy` and commit
   - Tag already exists: Use a different version number
   - Tests failing: Fix test failures before releasing

### Coverage Too Low

If coverage check fails:

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View detailed coverage
go tool cover -html=coverage.out

# Find uncovered code
go tool cover -func=coverage.out | grep -v 100.0%
```

Add tests for uncovered code and rerun.

## Post-Release Tasks

After successful release:

1. **Verify pkg.go.dev:**
   - Check: `https://pkg.go.dev/github.com/rcommerz/logger-go`
   - Verify documentation renders correctly
   - Check examples are visible

2. **Update documentation:**
   - Update installation badges if needed
   - Announce release in relevant channels
   - Update dependent projects

3. **Monitor issues:**
   - Watch for bug reports
   - Respond to questions
   - Plan next release

## GitHub Actions Workflows

### test.yml
- **Trigger**: Push to main/master/develop, Pull requests
- **Purpose**: Run tests across multiple Go versions
- **Checks**:
  - Tests pass on Go 1.21, 1.22, 1.23
  - Code formatting (gofmt)
  - go vet passes
  - Coverage ≥95%
  - golangci-lint passes

### release.yml
- **Trigger**: Push tags (v*.*.*), Manual dispatch
- **Purpose**: Create GitHub releases and publish package
- **Steps**:
  1. Validate version format
  2. Run tests with coverage
  3. Check required files
  4. Validate go.mod
  5. Create/push tag (if dispatch)
  6. Create GitHub Release
  7. Update on pkg.go.dev

### release-drafter.yml
- **Trigger**: Push to main/master, Pull requests
- **Purpose**: Auto-draft release notes from PRs
- **Features**:
  - Categorizes changes by labels
  - Auto-increments version
  - Generates changelog

## Maintenance

### Updating Dependencies

```bash
# Check for updates
go list -u -m all

# Update specific dependency
go get github.com/gofiber/fiber/v2@latest

# Update all dependencies
go get -u ./...

# Tidy and verify
go mod tidy
go mod verify

# Run tests
go test ./...
```

### Security Updates

Monitor for security advisories:

```bash
# Check for vulnerabilities
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

If vulnerabilities found, update dependencies and release a patch version.

## Support

- **Issues**: https://github.com/rcommerz/logger-go/issues
- **Discussions**: https://github.com/rcommerz/logger-go/discussions
- **Documentation**: https://pkg.go.dev/github.com/rcommerz/logger-go

## Additional Resources

- [Go Modules Reference](https://go.dev/ref/mod)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [pkg.go.dev About](https://pkg.go.dev/about)

git add .
git commit -m "feat: add new feature"

# Create new tag
git tag -a v1.1.0 -m "Release v1.1.0 - New features"

# Push changes and tag
git push origin master
git push origin v1.1.0
```

## Verification

After publishing, verify the package is accessible:

```bash
# In a new directory
mkdir test-logger
cd test-logger
go mod init test

# Install the package
go get github.com/rcommerz/logger-go@v1.0.0

# Verify it appears in go.mod
cat go.mod
```

## Package Discovery

Your package will be discoverable at:

- **pkg.go.dev**: <https://pkg.go.dev/github.com/rcommerz/logger-go>
- **Go Packages**: <https://go-packages.org/github.com/rcommerz/logger-go>
- **GitHub**: <https://github.com/rcommerz/logger-go>

## Maintenance

### Update Documentation

Keep these files updated:

- `README.md` - Usage and examples
- `CHANGELOG.md` - Version history
- `CONTRIBUTING.md` - Contribution guidelines
- `LICENSE` - License information

### Handling Issues

1. Monitor GitHub Issues
2. Respond to user questions
3. Accept pull requests
4. Release patches for critical bugs

### Security Updates

For security vulnerabilities:

1. Fix the issue
2. Create a patch release
3. Update CHANGELOG.md
4. Create GitHub Security Advisory

## Common Commands

```bash
# List all tags
git tag -l

# Delete a tag locally
git tag -d v1.0.0

# Delete a tag remotely
git push origin :refs/tags/v1.0.0

# View module info
go list -m github.com/rcommerz/logger-go@latest

# Check available versions
go list -m -versions github.com/rcommerz/logger-go
```

## Troubleshooting

### Package Not Found

If users report "package not found":

1. Verify tag exists: `git tag -l`
2. Verify pushed: `git ls-remote --tags`
3. Check module path in go.mod
4. Wait 10-15 minutes for proxy cache

### Import Path Issues

Ensure:

- Module path matches repository: `github.com/rcommerz/logger-go`
- No capital letters in repository name
- Repository is public

### Documentation Not Updating

Force pkg.go.dev refresh:

1. Visit: `https://pkg.go.dev/github.com/rcommerz/logger-go?tab=versions`
2. Click on your version
3. Wait for indexing to complete

## Best Practices

1. ✅ Always run tests before tagging
2. ✅ Update CHANGELOG.md for each release
3. ✅ Use semantic versioning
4. ✅ Write clear commit messages
5. ✅ Keep backward compatibility
6. ✅ Document breaking changes
7. ✅ Respond to issues promptly
8. ✅ Accept community contributions

## Resources

- [Go Modules Documentation](https://go.dev/doc/modules)
- [Publishing Go Modules](https://go.dev/blog/publishing-go-modules)
- [Semantic Versioning](https://semver.org/)
- [pkg.go.dev](https://pkg.go.dev/)
