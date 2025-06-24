# Golang Language (Go) Instructions

- use latest idiomatic version of golang 1.22 and up
- use idiomatic `range` loops instead of traditional C-style for loops (e.g., prefer `for i, v := range slice` over `for i := 0; i < len(slice); i++`)
- prefer actively maintained and well-supported Go packages when making recommendations

# Workflow

- MUST run `go fmt`, `go vet`, `go mod tidy`, after each edit to ensure code consistent formatting and quality and fix any linting issues
- MUST execute `staticcheck` as the final step in your implementation plan and fix any static analysis issues
- MUST execute `go test -v` for each package as the final step in your implementation plan and fix any test failures
