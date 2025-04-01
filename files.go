package main

import (
	"database/sql"
	"log"
)

type File struct {
	Id    int64
	dirId int64
}

func insertFile(dirId int64, fp string) (fileId int64) {
	result, err := tx.Exec(`INSERT INTO files (dir_id, filepath, updated_at) VALUES (?, ?, DATETIME('now'))`, dirId, fp)
	if err != nil {
		log.Fatalf("Error inserting file (dir_id=%d, filepath='%s'): %v\n", dirId, fp, err)
	}
	fileId, err = result.LastInsertId()
	if err != nil {
		log.Fatalf("Error accessing file 'id' field (dir_id=%d, filepath='%s'): %v\n", dirId, fp, err)
	}
	return fileId
}

func updateFileTimestamp(dirId int64, fp string) (fileId int64) {
	_, err := tx.Exec(`UPDATE files SET updated_at=DATETIME('now') WHERE dir_id=? AND filepath=?`, dirId, fp)
	if err != nil {
		log.Fatalf("Error updating file (dir_id=%d, filepath='%s'): %v\n", dirId, fp, err)
	}
	return selectFileFromDirIdAndFilepath(dirId, fp)
}

func maybeInsertFile(dirId int64, relPath string) (file File) {
	fileId := selectFileFromDirIdAndFilepath(dirId, relPath)
	if fileId != 0 {
		// We already got one!
		goto end
	}
	fileId = insertFile(dirId, relPath)
end:
	return File{
		dirId: dirId,
		Id:    fileId,
	}
}

func getFileId(dirId int64, fp string) (fileId int64, err error) {
	var stmt *sql.Stmt
	stmt, err = getTx().Prepare(`SELECT id FROM files WHERE dir_id = ? AND filepath = ?`)
	defer mustClose(stmt)
	if err != nil {
		log.Fatalf("Error preparing getFileId statement: %v\n", err)
	}
	err = stmt.QueryRow(dirId, fp).Scan(&fileId)
	if err != nil {
		log.Fatalf("Error querying fp %s in directort %d: %v\n", fp, dirId, err)
	}
	return fileId, err
}

func selectFileFromDirIdAndFilepath(dirId int64, fp string) (fileId int64) {
	_ = tx.QueryRow("SELECT id FROM files WHERE dir_id =? AND filepath = ?", dirId, fp).Scan(&fileId)
	// TODO Inspect sentinel error here to handle fail on non-expected errors
	return fileId
}

func selectFilesCountFromDirId(dirId int64) (count int) {
	// Check if this directory has files
	err := tx.QueryRow("SELECT COUNT(*) FROM files WHERE dir_id = ?", dirId).Scan(&count)
	if err != nil {
		log.Fatalf("Error selecting groups (dir_id=%d): %v", dirId, err)
	}
	return count
}

func deleteFilesForDirId(dirId int64) {
	_, err := tx.Exec("DELETE FROM files WHERE dir_id = ?", dirId)
	if err != nil {
		log.Fatalf("Error clearing files (dir_id=%d): %v\n", dirId, err)
	}
}

func deleteFilesForDirIdOlderThan(dirId int64, timestamp string) {
	_, err := tx.Exec(`DELETE FROM files WHERE dir_id=? AND updated_at < ?`, dirId, timestamp)
	if err != nil {
		log.Fatalf("Error deleting files (dir_id=%d, updated_id<'%s'): %v\n", dirId, timestamp, err)
	}
}
