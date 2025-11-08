# Go Best Practices for AI Agents

This document provides guidelines for AI agents working on Go code in this repository.

## Code Formatting and Style

### Use gofmt and goimports
- Always format code using `gofmt` or `go fmt`
- Use `goimports` to automatically manage import statements
- Run these tools before committing any Go code

### Follow Standard Go Style
- Use the official [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) as a guide
- Follow the [Effective Go](https://go.dev/doc/effective_go) guidelines
- Keep line length reasonable (typically under 100-120 characters)

## Naming Conventions

### Package Names
- Use short, lowercase, single-word names
- Avoid underscores, hyphens, or mixedCaps
- Use common abbreviations when appropriate (e.g., `http`, `fmt`, `io`)

### Variable and Function Names
- Use mixedCaps or MixedCaps (camelCase or PascalCase)
- Exported names must start with a capital letter
- Unexported names must start with a lowercase letter
- Use short names for short-lived variables (e.g., `i`, `j` for loop counters)
- Use descriptive names for longer-lived variables and exported functions

### Acronyms
- Use all caps for acronyms (e.g., `URL`, `HTTP`, `API`)
- Examples: `ServeHTTP`, `parseURL`, `IDProcessor`

## Error Handling

### Always Check Errors
```go
// GOOD
f, err := os.Open("file.txt")
if err != nil {
    return fmt.Errorf("failed to open file: %w", err)
}
defer f.Close()

// BAD - don't ignore errors
f, _ := os.Open("file.txt")
```

### Use Error Wrapping
- Wrap errors with context using `fmt.Errorf` with `%w` verb
- This preserves the error chain for `errors.Is` and `errors.As`

### Create Custom Errors When Appropriate
- Define sentinel errors as package-level variables
- Use custom error types for complex error scenarios

## Package Organization

### Keep Packages Focused
- Each package should have a single, clear purpose
- Avoid circular dependencies between packages
- Use internal packages for code that shouldn't be imported externally

### Main Package
- Keep the `main` package minimal
- Move business logic into separate packages
- The `main` function should primarily handle initialization and orchestration

## Testing

### Write Table-Driven Tests
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   int
        want    int
        wantErr bool
    }{
        {"positive", 5, 10, false},
        {"zero", 0, 0, false},
        {"negative", -5, 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Function() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Function() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test File Organization
- Name test files with `_test.go` suffix
- Place tests in the same package as the code being tested
- Use `package <name>_test` for black-box testing

### Use Testing Helpers
- Use `t.Helper()` in test helper functions
- Use `t.Cleanup()` for cleanup operations
- Use subtests with `t.Run()` for better organization

## Documentation

### Document Exported Identifiers
- Every exported function, type, constant, and variable should have a doc comment
- Doc comments should be complete sentences starting with the name
- Example: `// Open opens the named file for reading.`

### Package Documentation
- Every package should have a package comment
- Place it in a file (commonly `doc.go` or the main file)
- Start with "Package <name>" followed by a description

### Use Examples
- Provide testable examples for complex functionality
- Place examples in `example_test.go` files

## Concurrency

### Use Channels for Communication
- Prefer channels over shared memory for goroutine communication
- Use buffered channels carefully and document the buffer size rationale

### Protect Shared State
- Use `sync.Mutex` or `sync.RWMutex` to protect shared state
- Keep critical sections small
- Consider using `sync.Once` for one-time initialization

### Avoid Goroutine Leaks
- Always ensure goroutines can exit
- Use context for cancellation
- Use `defer` to ensure cleanup happens

### Context Usage
```go
// GOOD - respect context cancellation
func processData(ctx context.Context, data []byte) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // process data
    }
    return nil
}
```

## Performance

### Preallocate Slices When Size is Known
```go
// GOOD
items := make([]Item, 0, expectedSize)

// LESS EFFICIENT
items := []Item{}
```

### Use String Builder for Concatenation
```go
// GOOD
var sb strings.Builder
for _, s := range strings {
    sb.WriteString(s)
}
result := sb.String()

// BAD - creates many temporary strings
result := ""
for _, s := range strings {
    result += s
}
```

### Avoid Unnecessary Allocations
- Reuse buffers when possible
- Use `sync.Pool` for frequently allocated objects
- Profile before optimizing

## Code Organization

### Struct Field Order
- Group related fields together
- Consider ordering by size for better memory alignment (largest first)
- Place exported fields before unexported fields (optional convention)

### Function Order
- Place exported functions before unexported functions
- Group related functions together
- Place helper functions near where they're used

### Interfaces
- Keep interfaces small and focused
- Define interfaces where they're used, not where they're implemented
- Prefer many small interfaces over one large interface

## Build and Dependencies

### Use Go Modules
- Keep `go.mod` and `go.sum` up to date
- Run `go mod tidy` to clean up dependencies
- Use `go mod vendor` if vendoring dependencies

### Minimize Dependencies
- Avoid adding dependencies for trivial functionality
- Prefer standard library when possible
- Keep dependencies up to date for security

## Security

### Input Validation
- Validate all external input
- Use parameterized queries for database operations
- Sanitize input before logging

### Avoid Common Pitfalls
- Don't use `unsafe` package unless absolutely necessary
- Be careful with `reflect` package - it can bypass type safety
- Use `crypto/rand` for cryptographic operations, not `math/rand`

## Linting and Static Analysis

### Use Standard Tools
- Run `go vet` to catch common mistakes
- Use `golangci-lint` for comprehensive linting
- Run `go test -race` to detect race conditions
- Use `staticcheck` for additional static analysis

### Address All Warnings
- Fix or explicitly ignore all linter warnings
- Document why warnings are ignored with `//nolint` comments

## Version Compatibility

- Maintain compatibility with the Go version specified in `go.mod`
- Test with the minimum supported Go version
- Use build tags for version-specific code if needed

## Before Committing

Run this checklist:
1. [ ] Code is formatted with `gofmt`
2. [ ] Imports are organized with `goimports`
3. [ ] `go vet` passes without errors
4. [ ] All tests pass: `go test ./...`
5. [ ] Race detector passes: `go test -race ./...`
6. [ ] Documentation is updated
7. [ ] Error handling is comprehensive
8. [ ] No unnecessary dependencies added

## Additional Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Go Proverbs](https://go-proverbs.github.io/)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
