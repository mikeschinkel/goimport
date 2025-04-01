# goimports

A quick-and-dirty CLI tool to analyze Go imports across multiple repositories.

## What It Does

This tool scans Go files in a directory, tracks import statements, and stores the data in a SQLite database located at `~/.config/goimports/goimports.db`. It helps identify:

- Most frequently used imports across projects
- Import patterns within repositories
- Dependencies across your codebase

## Installation

```
go install github.com/mikeschinkel/goimports@latest
```

## Usage

The `-mode` parameter is required and specifies what output to display.

```
# Display imports in current directory files
goimports -mode=files

# Display files for each import in current directory
goimports -mode=imports

# Display imports grouped by directory
goimports -mode=dirs

# Display imports grouped by repository
goimports -mode=repos

# Only update the database, no output
goimports -mode=none

# Scan specific directory and show files with imports
goimports -dir=/path/to/go/project -mode=files

# Enable verbose output
goimports -v -mode=files
```

Run `goimports` without arguments to see all available options.

## Bulk Processing
The following Bash command will run `goimports` in each child subdirectory of the current directory allowing you to slurp up import information on multiple repos if you have them all cloned to a single directory:

```shell
find . \
   -maxdepth 1 \
   -type d \
   -not -path "." \
   | xargs \
      -I{} \
      sh -c 'cd "{}" && echo "Processing $(basename {})" && goimports -mode none'
```

## Note on Code Quality

⚠️ **Disclaimer**: This project was intentionally built as a quick utility tool for personal use and deliberately avoids following certain Go best practices:

- Uses `log.Fatal()` instead of proper error handling in many places
- Minimal test coverage
- Limited input validation
- Takes some shortcuts with SQL queries

This does not represent my professional coding standards - just a practical tool I needed quickly. Please don't judge my Go skills based on this repository.

## Database

The SQLite database is stored at `~/.config/goimports/goimports.db` and includes tables for:
- Directories
- Files
- Imports
- File imports

## License

MIT