# Golang Language (Go) Instructions

- use latest idiomatic version of golang 1.22 and up
- use idiomatic `range` loops instead of traditional C-style for loops (e.g., prefer `for i, v := range slice` over `for i := 0; i < len(slice); i++`)
- prefer actively maintained and well-supported Go packages when making recommendations

# Workflow

- MUST run `go fmt` and `go vet`, `go test -v` after each edit to ensure code consistent formatting and quality
- MUST execute `go test -v` for each package as the final step in your implementation plan
