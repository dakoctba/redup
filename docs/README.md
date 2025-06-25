# Redup - Developer Documentation

This documentation contains technical information about the development, compilation, and source code of Redup.

## Development Requirements

- Go 1.21 or higher
- Git
- [GoReleaser](https://goreleaser.com/install/) (for releases)

## Installation for Development

### Prerequisites

- Go 1.21 or higher

### Cloning the Repository

```bash
git clone https://github.com/dakoctba/redup.git
cd redup
```

### Installing Dependencies

```bash
go mod tidy
```

## Compilation

### Using Make

The project includes a Makefile to facilitate the development process:

```bash
# Build the project
make build

# Run tests
make test

# Run the application with arguments
make run ARGS="--dry-run ."

# Clean build files
make clean

# Install locally
make install

# Uninstall
make uninstall

# Display Makefile help
make help
```

### Using Go Directly

```bash
# Build the project
go build

# Install the binary
go install
```

## Project Structure

```
.
├── cmd/
│   └── root.go      # Main command and configuration using Cobra
├── pkg/
│   ├── scanner.go    # File scanning logic
│   ├── hasher.go     # Checksum calculation
│   ├── deduplicator.go # Duplicate detection
│   ├── backup.go     # Backup management
│   ├── menu.go       # Interactive menu
│   ├── reporter.go   # Statistics and reporting
│   ├── gitignore.go  # .gitignore processing
│   └── config.go     # Configuration management
├── bin/
│   ├── release.sh         # Release creation script
│   └── update_version.sh  # Version update script
├── docs/
│   └── README.md          # Technical documentation
├── main.go          # Entry point
├── go.mod           # Dependency management
├── .goreleaser.yml  # GoReleaser configuration
├── Makefile         # Task automation
└── README.md        # User documentation
```

## Running Tests

```bash
go test ./pkg/...
```

## Version Management and Releases

Redup uses Git tags for version control following Semantic Versioning (SemVer).

### Release Process

The project includes an interactive process to create new versions and releases:

```bash
make release
```

This command will:
1. Show the current version (based on Git tags)
2. Suggest options for the next version following Semantic Versioning
3. Allow you to choose between patch, minor, major, or a custom version
4. Request confirmation of the operation
5. Create a Git tag
6. Push the tag to GitHub
7. Run GoReleaser to generate the release

### Release Scripts

#### update_version.sh

This script:
- Gets the latest version from Git tags
- Suggests the next version (patch, minor, or major)
- Allows manual version selection
- Creates a Git tag and publishes it to the remote repository
- Calls GoReleaser to create the release

#### release.sh

This script:
- Manages the release process
- Can be used directly (without interactivity)
- Accepts options like `--snapshot`, `--clean`, `--no-clean`
- Executes GoReleaser with the appropriate settings

```bash
# Run directly (uses Git tag version)
bin/release.sh

# Create a snapshot (for testing)
bin/release.sh --snapshot

# Preserve the dist/ directory
bin/release.sh --no-clean

# Use a specific version
bin/release.sh --version=1.2.3
```

### GoReleaser

Redup uses [GoReleaser](https://goreleaser.com) to automate the release process for multiple platforms.

Requirements:
- GoReleaser installed
- GitHub token configured (in the .env file)

The GoReleaser configuration is in the `.goreleaser.yml` file.

### Manual Release

If you prefer to run GoReleaser directly:

```bash
# Export GitHub token
export GITHUB_TOKEN=your_github_token

# Run GoReleaser
goreleaser release --clean
```

## Environment Configuration

Redup uses environment variables for configuration, which can be defined in a `.env` file:

```
GITHUB_TOKEN=your_github_token_here
DEBUG=true
LOG_LEVEL=info
```

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add Amazing Feature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Implementation Details

### File Processing

The `pkg` package contains the main logic for file processing:
- Recursive directory listing with `.gitignore` support
- Checksum calculation (SHA-256 and MD5)
- Duplicate detection and grouping
- Safe backup operations
- Interactive menu system

### Command Line Interface

The CLI is implemented using the [Cobra](https://github.com/spf13/cobra) library, with commands defined in `cmd/root.go`.

### Key Features

1. **Scanner**: Recursively scans directories, respecting `.gitignore` rules
2. **Hasher**: Calculates checksums for duplicate detection
3. **Deduplicator**: Groups files by content using checksums
4. **Backup Manager**: Safely moves duplicate files to timestamped directories
5. **Interactive Menu**: User-friendly interface for managing duplicates
6. **Reporter**: Generates statistics and JSON output

### Dependencies

- **Cobra**: CLI framework
- **Godotenv**: Environment variable loading
- **Standard library**: File operations, crypto, etc.

## Testing Strategy

The project includes comprehensive tests for:
- File scanning functionality
- Checksum calculation
- Duplicate detection
- `.gitignore` processing
- Backup operations
- Interactive menu

Run tests with:
```bash
make test
```

## Code Style

The project follows Go conventions and best practices:
- Use `gofmt` for code formatting
- Follow Go naming conventions
- Include comprehensive comments
- Write unit tests for all functions
- Use meaningful variable and function names
