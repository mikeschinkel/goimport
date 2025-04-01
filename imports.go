package main

import (
	"cmp"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type importInfos []*importInfo

func (ii importInfos) sort(st SortType) {
	switch st {
	case ImportSort:
		slices.SortFunc(ii, func(a, b *importInfo) int {
			return strings.Compare(a.path, b.path)
		})
	case CountSort:
		fallthrough
	default:
		slices.SortFunc(ii, func(a, b *importInfo) int {
			order := -cmp.Compare(a.count, b.count)
			if order == 0 {
				order = strings.Compare(a.path, b.path)
			}
			return order
		})
	}
}

type importInfo struct {
	path  string
	count int
}

func updateImportsForIds(ids Ids) {
	_, err := tx.Exec(
		"UPDATE dirs SET repo_id = ?, updated_at = DATETIME('now') WHERE id = ? AND repo_id <> ?", ids.repoId, ids.dirId, ids.repoId,
	)
	if err != nil {
		log.Fatalf("Error updating directory metadata: %v\n", err)
	}
}

func selectImportByImportPath(importPath string) (importId int64) {
	// Check if this directory has files
	err := tx.QueryRow(`SELECT id FROM imports WHERE import_path = ?`, importPath).Scan(&importId)
	if err != nil {
		log.Fatalf("Error selecting import (import_path='%s'): %v", importPath, err)
	}
	return importId
}

func verifyImportByIdentity(importId int64) {
	// Check if this directory has files
	err := tx.QueryRow(`SELECT id FROM imports WHERE id = ?`, importId).Scan(&importId)
	if err != nil {
		log.Fatalf("Error selecting import (import_id=%d): %v", importId, err)
	}
}

func insertImport(importPath string) (repoId int64) {
	result, err := tx.Exec(`INSERT INTO imports (import_path, updated_at) VALUES (?, DATETIME('now'))`, importPath)
	if err != nil {
		// Probably already exists
		// TODO Inspect sentinel error here to handle fail on non-expected errors
		goto end
	}
	repoId, err = result.LastInsertId()
	if err != nil {
		log.Fatalf("Error getting import ID: %v\n", err)
	}
end:
	return repoId
}

func maybeInsertImports(ids Ids) {
	// Walk the directory tree and collect imports
	processedFiles := 0
	err := filepath.Walk(inputs.rootPath, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			// We got an error on entry? Bail...
			goto end
		}
		if file.IsDir() {
			// No need to read anything from a directory
			goto end
		}
		if !strings.HasSuffix(path, ".go") {
			// No need to read non-Go files
			goto end
		}
		insertOrUpdateFileImports(path, ids)
		processedFiles++
		if processedFiles%100 == 0 {
			logf("Processed %d files\n", processedFiles)
		}
	end:
		return err
	})
	if err != nil {
		log.Fatalf("Error walking directory: %v\n", err)
	}
	logf("Database updated with imports from %d files\n", processedFiles)
	return
}
