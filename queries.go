package main

import (
	"database/sql"
	"fmt"
	"log"
)

func queryAllImports(sort SortType) (*sql.Rows, error) {
	var sortExpr string
	switch sort {
	case FileSort:
		sortExpr = "import_count DESC,import_path"
	case ImportSort:
		sortExpr = "import_path"
	default:
		log.Fatalf("Invalid sort type for '%s' mode: %s", AllOutput, sort)
	}
	return db.Query(fmt.Sprintf(`
		SELECT import_path,COUNT(import_path) AS import_count
		FROM repo_imports
		GROUP BY import_path
		ORDER BY %s;
	`, sortExpr))
}

func queryFileImportsForDirId(dirId int64) (*sql.Rows, error) {
	return db.Query(`
		SELECT f.filepath, GROUP_CONCAT(i.import_path, CHAR(10)) AS imports
		FROM files f
		JOIN file_imports fi ON f.id = fi.file_id
		JOIN imports i ON fi.import_id = i.id
		WHERE f.dir_id = ?
		GROUP BY f.filepath
		ORDER BY f.filepath
	`, dirId)
}

func queryImportFilesForDirId(dirId int64) (*sql.Rows, error) {
	return db.Query(`
		SELECT 
			i.import_path,
			GROUP_CONCAT(f.filepath, CHAR(10)) AS filepaths
		FROM imports i
		JOIN file_imports fi ON i.id = fi.import_id
		JOIN files f ON fi.file_id = f.id
		WHERE f.dir_id = ?
		GROUP BY i.import_path
		ORDER BY i.import_path
	`, dirId)
}
