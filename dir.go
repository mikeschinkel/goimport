package main

import (
	"database/sql"
	"log"
)

func selectDirFromLocalPath(localPath string) (dirId int64) {
	_ = tx.QueryRow("SELECT id FROM dirs WHERE local_path = ?", localPath).Scan(&dirId)
	// TODO Inspect sentinel error here to handle fail on non-expected errors
	return dirId
}

func insertDirFromInputs(repoId int64) (dirId int64) {
	// Insert new directory
	result, err := tx.Exec(
		"INSERT INTO dirs (repo_id, local_path, updated_at) VALUES (?, ?, DATETIME('now'))",
		repoId,
		inputs.rootPath,
	)
	if err != nil {
		log.Fatalf("Error inserting directory data: %v\n", err)
	}

	dirId, err = result.LastInsertId()
	if err != nil {
		log.Fatalf("Error getting directory ID: %v\n", err)
	}
	return dirId
}

func getDirId(repoId int64) (dirId int64) {
	dirId = selectDirFromLocalPath(inputs.rootPath)
	if dirId == 0 {
		dirId = insertDirFromInputs(repoId)
		goto end
	}
end:
	return dirId
}

func displayDirImports() {
	rows, err := queryDirImports()
	if err != nil {
		log.Fatalf("Error querying database: %v\n", err)
	}
	defer mustClose(rows)
	dirs := newImportGroups()
	dirs.collect(rows)
	dirs.display("dir")
}

func queryDirImports() (*sql.Rows, error) {
	return db.Query(`
   SELECT
      d.local_path AS dir,
      i.import_path AS import,
      COUNT(i.id) AS import_count
   FROM dirs d
		 JOIN files f ON d.id = f.dir_id
		 JOIN file_imports fi ON f.id = fi.file_id
		 JOIN imports i ON fi.import_id = i.id
   GROUP BY d.local_path,i.import_path
   ORDER BY d.local_path,i.import_path
	`)
}
