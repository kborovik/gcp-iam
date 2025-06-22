# Golang Language (Go) Instructions

- use latest idiomatic version of golang 1.22 and up
- use idiomatic `range` loops instead of traditional C-style for loops (e.g., prefer `for i, v := range slice` over `for i := 0; i < len(slice); i++`)
- when you only need the value, use `for _, v := range slice`
- when you only need the index, use `for i := range slice`
- MUST run `go fmt` and `go vet` after each edit to ensure code consistent formatting and quality
- run advanced Go linter `staticcheck` for additional static analysis
- prefer actively maintained and well-supported Go packages when making recommendations

# Application CLI Instructions

- test application by running `go run main.go` and follow help prompt
