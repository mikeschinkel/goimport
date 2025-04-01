package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

func deleteFileImportsForDirId(dirId int64) {
	_, err := tx.Exec(`DELETE FROM file_imports WHERE file_id IN (SELECT id FROM files WHERE dir_id = ?)`, dirId)
	if err != nil {
		log.Fatalf("Error deleting file imports (dir_id=%d): %v\n", dirId, err)
	}
}

func deleteFileImportsForDirIdOlderThan(dirId int64, timestamp string) {
	_, err := tx.Exec(`DELETE FROM file_imports WHERE file_id IN (SELECT id FROM files WHERE dir_id=? AND updated_at < ?)`, dirId, timestamp)
	if err != nil {
		log.Fatalf("Error deleting file imports (dir_id=%d, updated_id<'%s'): %v\n", dirId, timestamp, err)
	}
}

func selectFileImportForFileIdAndImportId(fileId, importId int64) (fileImportId int64) {
	_ = tx.QueryRow(`SELECT id FROM file_imports WHERE file_id = ? AND import_id = ?`, fileId, importId).Scan(&fileImportId)
	return fileImportId
}

func insertFileImport(fileId, importId int64) (fileImportId int64) {
	// Insert new repository
	result, err := tx.Exec(`INSERT INTO file_imports (file_id, import_id, updated_at) VALUES (?,?, DATETIME('now'))`, fileId, importId)
	if err != nil {
		// Probably already exists
		// TODO Inspect sentinel error here to handle fail on non-expected errors
		goto end
	}
	fileImportId, err = result.LastInsertId()
	if err != nil {
		log.Fatalf("Error inserting file import (file_id=%d, import_id=%d): %v\n", fileId, importId, err)
	}
end:
	return fileImportId
}

func updateFileImportTimestamp(fileImportId int64) {
	// Insert new repository
	_, err := tx.Exec(`UPDATE file_imports SET updated_at = DATETIME('now') WHERE id = ?`, fileImportId)
	if err != nil {
		log.Fatalf("Error updating file import (file_import_id=%d): %v\n", fileImportId, err)
	}
}

func insertOrUpdateFileImports(path string, dirId int64) {
	var err error
	var fileSet *token.FileSet
	var goFileAST *ast.File
	var relPath string
	var file File

	fileSet = token.NewFileSet()
	goFileAST, err = parser.ParseFile(fileSet, path, nil, parser.ImportsOnly)
	if err != nil {
		log.Fatalf("Error parsing %s: %v\n", path, err)
	}

	// Get relative path for storing
	relPath, err = filepath.Rel(inputs.rootPath, path)
	if err != nil {
		relPath = path
	}

	// Insert ast record
	file = maybeInsertFile(dirId, relPath)
	if file.Id != 0 {
		file = File{
			Id:    updateFileTimestamp(dirId, relPath),
			dirId: dirId,
		}
	}
	for _, imp := range goFileAST.Imports {
		fileImportId := selectFileFromDirIdAndFilepath(file.Id, imp.Path.Value)
		if fileImportId == 0 {
			maybeInsertFileImport(file.Id, imp.Path.Value)
			continue
		}
		updateFileImportTimestamp(fileImportId)
	}
}

func maybeInsertFileImport(fileId int64, imp string) int64 {
	importPath := strings.Trim(imp, `"`)
	// Insert or get import record
	importId := insertImport(importPath)
	if importId == 0 {
		importId = selectImportByImportPath(importPath)
	} else {
		verifyImportByIdentity(importId)
	}
	return insertFileImport(fileId, importId)
}
