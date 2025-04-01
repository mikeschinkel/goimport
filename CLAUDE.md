# CLAUDE.md - goimports project guide

## Build / Test / Lint Commands
- Build: `go build`
- Test: `go test ./...`
- Test single file: `go test -run TestName`
- Format code: `go fmt ./...`
- Reset database: `./resetdb.sh`

## Code Style Guidelines
- Use Go 1.24+ features like `cmp` package for comparisons
- Import order: standard library first, then third-party packages
- Error handling: early returns with `goto end` pattern for cleanup
- Naming: camelCase for variables, PascalCase for exported functions
- Database: Use prepared statements with SQLite and `sql.Tx` transactions
- Comments: Add comments for exported functions
- Interfaces: Small, focused interfaces (see `rollbacker`, `Committer`)
- Error fatals: Use `log.Fatalf` for unrecoverable errors
- SQL paths: SQL scripts stored in `/migrations` directory
- Paths: Use `filepath.Join` for platform-compatible paths