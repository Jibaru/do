.PHONY: run test coverage fmt deps

# Variables
MAIN_PACKAGE := ./cmd/do/main.go
TEST_COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html

# Run the application
run:
	go run $(MAIN_PACKAGE)

# Run all tests with verbose output
test:
	go test -v ./...

# Run tests with coverage and generate coverage report in HTML
coverage:
	go test -v -coverprofile=$(TEST_COVERAGE_OUT) ./...
	go tool cover -html=$(TEST_COVERAGE_OUT) -o $(COVERAGE_HTML)

# Format the code
fmt:
	go fmt ./...

# Install dependencies
deps:
	go mod tidy
	go mod download
