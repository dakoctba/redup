.PHONY: build test clean install

# Build the application
build:
	go build -o redup .

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Clean build artifacts
clean:
	rm -f redup
	rm -f *.exe
	rm -f *.dll
	rm -f *.so
	rm -f *.dylib

# Install the application
install:
	go install .

# Build for release
release:
	GOOS=linux GOARCH=amd64 go build -o redup-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o redup-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o redup-windows-amd64.exe .

# Default target
all: build test
