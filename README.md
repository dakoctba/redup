# Redup - Duplicate File Manager

Redup is a command-line tool that allows you to find and manage duplicate files by content using checksums (SHA-256 or MD5), respecting .gitignore rules and providing safe backup options.

## Features

- Recursive directory scanning
- Content-based duplicate detection using SHA-256 or MD5 checksums
- Automatic .gitignore support
- Minimum file size filtering
- Safe backup system with timestamped directories
- Dry-run mode for simulation
- JSON output format
- Interactive duplicate management
- Detailed processing statistics

## Installation

### Installation via go install

```bash
go install github.com/dakoctba/redup@latest
```

### Download Binaries

Visit the [releases page](https://github.com/dakoctba/redup/releases) to download the latest version compiled for your operating system.

## Usage

```bash
redup [options] [directory]
```

### Options

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--dir` | `-d` | Directory to scan (default: current working directory) | `--dir ~/Documents` |
| `--checksum` | `-c` | Checksum algorithm (sha256\|md5) | `--checksum md5` |
| `--min-size` | `-s` | Minimum file size to consider in bytes | `--min-size 1048576` |
| `--backup-dir` | `-b` | Base directory for backup | `--backup-dir ~/backups` |
| `--dry-run` | `-n` | Simulate actions without moving files | `--dry-run` |
| `--json` | `-j` | Output results in JSON format | `--json` |
| `--version` | `-v` | Show version number | `--version` |

### Commands

| Command | Description | Example |
|---------|-------------|---------|
| `completion` | Generate autocompletion script | `redup completion bash` |
| `version` | Display detailed application version | `redup version` |

### Examples

```bash
# Scan current directory for duplicates
redup

# Scan specific directory with minimum size filter
redup --dir ~/Documents --min-size 1048576

# Use MD5 checksum with dry-run simulation
redup --checksum md5 --dry-run ~/Pictures

# Export results to JSON file
redup --json ~/Music > duplicates.json

# Use custom backup directory
redup --backup-dir ~/backups ~/Downloads

# Show version information
redup -v
# or
redup version
```

## Output Behavior

Redup has different output behaviors depending on how it's used:

1. **Default mode** (`redup`):
   - Shows help information
   - No scanning performed

2. **Scanning mode** (`redup [directory]`):
   - Scans directory for duplicates
   - Shows summary of found duplicates
   - If duplicates found, starts interactive mode for management

3. **JSON mode** (`redup --json [directory]`):
   - Outputs detailed JSON with all duplicate information
   - Suitable for programmatic processing

4. **Dry-run mode** (`redup --dry-run [directory]`):
   - Simulates the scanning process
   - Shows what would be done without making changes

## Gitignore Support

Redup automatically reads and respects the `.gitignore` file in your project directory. This means that files and directories listed in your `.gitignore` will be automatically excluded from scanning, including:

- `node_modules/`
- `.env` files
- Build artifacts
- Log files
- And any other patterns you've specified in your `.gitignore`

## Backup System

When duplicates are found and you choose to manage them, Redup creates a safe backup system:

1. **Timestamped backup directories**: Each backup session creates a new directory with timestamp
2. **Original file preservation**: Original files are moved to backup, not deleted
3. **Interactive selection**: You choose which files to keep and which to backup
4. **Safe operations**: All operations are reversible through the backup system

## Statistics

At the end of execution, Redup displays statistics about the processed files:

- Total number of files scanned
- Number of duplicate groups found
- Total size of duplicate files
- Number of files processed
- Backup operations performed

## Return Codes

| Code | Description |
|------|-------------|
| 0 | Successful execution |
| 1 | Usage error (invalid arguments) |
| 2 | Error scanning/processing files |

## Checksum Algorithms

Redup supports two checksum algorithms:

- **SHA-256** (default): More secure, slower for large files
- **MD5**: Faster, suitable for most duplicate detection scenarios

## Interactive Mode

When duplicates are found, Redup enters an interactive mode where you can:

- View detailed information about each duplicate group
- Select which files to keep
- Choose which files to backup
- Skip groups you don't want to process
- Exit the process at any time

## Developer Documentation

For information about development, compilation, and source code, see the [developer documentation](docs/README.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

- **Jackson** - [dakoctba](https://github.com/dakoctba)

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - Library for creating CLI applications in Go
- [Go](https://golang.org/) - Programming language
