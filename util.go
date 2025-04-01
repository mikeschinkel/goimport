package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func logf(format string, a ...interface{}) {
	if inputs.verbose {
		log.Printf(format, a...)
	}
}
func createInitialDatabase(db *sql.DB) {
	// Brand new database - create all tables
	query := mustLoadMigrationSQL(`migrations/version_1.sql`)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating schema: %v\n", err)
	}
}

type rollbacker interface {
	Rollback() error
}
type Committer interface {
	Commit() error
}

func mustRollback(r rollbacker) {
	err := r.Rollback()
	if errors.Is(err, sql.ErrTxDone) {
		goto end
	}
	if err != nil {
		log.Fatalf("Error attempting to rollback transaction: %v\n", err)
	}
end:
}
func mustCommit(c Committer) {
	err := c.Commit()
	if err != nil {
		log.Fatalf("Error attempting to commit transaction: %v\n", err)
	}
}
func mustClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatalf("Error closing connection: %v\n", err)
	}
}

func mustLoadSQL(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		// Fatal error if we can't load the SQL file
		log.Fatalf("Failed to load SQL file %s: %v", filename, err)
	}
	return string(content)
}

func getDBPath(dbPath string) string {
	var err error
	var homeDir, configDir string

	// Determine the database path
	if dbPath != "" {
		goto end
	}
	homeDir, err = os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v\n", err)
	}
	configDir = filepath.Join(homeDir, ".config", configSubdir)
	// Create config directory if it doesn't exist
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		log.Fatalf("Error creating config directory: %v\n", err)
	}
	dbPath = filepath.Join(configDir, defaultDBFile)
end:
	return dbPath
}

func findGitRoot(startPath string) (gr string, err error) {
	cmd := exec.Command("git", "-C", startPath, "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		goto end
	}
	gr = strings.TrimSpace(string(output))
end:
	return gr, err
}

func getGitOrigin(rootPath string) string {
	cmd := exec.Command("git", "-C", rootPath, "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	origin := strings.TrimSpace(string(output))

	// Normalize git origin URL
	if strings.HasPrefix(origin, "git@") {
		// Convert git@github.com:org/repo.git to github.com/org/repo
		re := regexp.MustCompile(`git@([^:]+):(.+)`)
		matches := re.FindStringSubmatch(origin)
		if len(matches) == 3 {
			domain := matches[1]
			path := strings.TrimSuffix(matches[2], ".git")
			origin = domain + "/" + path
		}
	} else if strings.HasPrefix(origin, "https://") {
		// Convert https://github.com/org/repo.git to github.com/org/repo
		origin = strings.TrimPrefix(origin, "https://")
		origin = strings.TrimSuffix(origin, ".git")
	}

	return origin
}

func getGoModule(rootPath string) string {
	goModPath := filepath.Join(rootPath, "go.mod")

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}

	return ""
}

func getWorkingDir(rootPath string) (wd string) {
	var err error
	if rootPath != "" && rootPath != "." {
		wd = rootPath
		goto end
	}
	wd, err = os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v\n", err)
	}
end:
	return wd
}

func dirExists(path string) (exists bool, err error) {
	var info os.FileInfo
	info, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			goto end
		}
		err = fmt.Errorf("error checking directory %s: %v", path, err)
		goto end
	}
	if !info.IsDir() {
		// Path exists but is not a directory
		goto end
	}
	exists = true
end:
	return exists, err
}
