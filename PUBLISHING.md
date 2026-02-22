# Publishing logger-go Package

This guide explains how to publish and maintain the logger-go package.

## Prerequisites

1. **GitHub Repository**
   - Repository should be public at `github.com/rcommerz/logger-go`
   - Repository must contain all source files

2. **Go Module**
   - `go.mod` file with correct module path
   - All dependencies properly declared

3. **Git Tags**
   - Semantic versioning (e.g., v1.0.0)
   - Annotated tags recommended

## Publishing Steps

### 1. Prepare the Release

Ensure all files are ready:

```bash
# Run tests
go test -v -cover

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy

# Verify module
go mod verify
```

### 2. Create Git Repository (First Time Only)

```bash
cd /home/miaki-maruf/workspace/personal/rcommerz/packages/logger-go

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

### 3. Create a Release Tag

```bash
# Create annotated tag
git tag -a v1.0.0 -m "Release v1.0.0 - Initial stable release"

# Push tag to GitHub
git push origin v1.0.0
```

### 4. Publish to pkg.go.dev

Once you push the tag to GitHub, pkg.go.dev will automatically index it. To trigger immediate indexing:

1. Visit: `https://pkg.go.dev/github.com/rcommerz/logger-go@v1.0.0`
2. Or use the Go proxy:

   ```bash
   GOPROXY=https://proxy.golang.org GO111MODULE=on go get github.com/rcommerz/logger-go@v1.0.0
   ```

### 5. Create GitHub Release

1. Go to: `https://github.com/rcommerz/logger-go/releases/new`
2. Select tag: `v1.0.0`
3. Release title: `v1.0.0 - Initial Release`
4. Description: Copy from CHANGELOG.md
5. Click "Publish release"

## Version Management

Follow [Semantic Versioning](https://semver.org/):

- **Major** (v2.0.0): Breaking changes
- **Minor** (v1.1.0): New features, backward compatible
- **Patch** (v1.0.1): Bug fixes, backward compatible

### Creating New Releases

```bash
# Make changes
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
