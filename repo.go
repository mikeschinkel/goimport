package main

import (
	"database/sql"
	"log"
)

func ensureRepoId() (repoId int64) {
	if inputs.gitOriginURL == "" {
		repoId = selectRepoIdFromModulePath(inputs.modulePath)
		goto end
	}
	repoId = selectRepoIdFromOriginURL(inputs.gitOriginURL)
	if repoId == 0 {
		repoId = insertRepoFromInputs()
		goto end
	}
	updateRepoModulePath(repoId, inputs.modulePath)
end:
	return repoId
}

func countRepos() (count int) {
	_ = db.QueryRow("SELECT COUNT(*) AS count FROM repos").Scan(&count)
	return count
}
func selectRepoIdFromOriginURL(origin string) (repoId int64) {
	_ = db.QueryRow("SELECT id FROM repos WHERE origin_url = ?", origin).Scan(&repoId)
	return repoId
}
func selectRepoIdFromModulePath(modulePath string) (repoId int64) {
	_ = db.QueryRow("SELECT id FROM repos WHERE origin_url IS NULL AND module_path = ?", modulePath).Scan(&repoId)
	return repoId
}

func insertRepoFromInputs() (repoId int64) {
	// Insert new repository
	result, err := tx.Exec(
		"INSERT INTO repos (origin_url, module_path, updated_at) VALUES (?, ?, DATETIME('now'))", inputs.gitOriginURL, inputs.modulePath,
	)
	if err != nil {
		log.Fatalf("Error inserting repository data: %v\n", err)
	}
	repoId, err = result.LastInsertId()
	if err != nil {
		log.Fatalf("Error getting repository ID: %v\n", err)
	}
	return repoId
}
func updateRepoModulePath(repoId int64, modulePath string) {
	// Insert new repository
	result, err := tx.Exec(
		"UPDATE repos SET module_path = ?, updated_at = DATETIME('now') WHERE id = ?", modulePath, repoId,
	)
	if err != nil {
		log.Fatalf("Error inserting repository data: %v\n", err)
	}

	repoId, err = result.LastInsertId()
	if err != nil {
		log.Fatalf("Error getting repository ID: %v\n", err)
	}
}

func queryRepoImports() (*sql.Rows, error) {
	return db.Query(`
	SELECT repo,import,COUNT(import) FROM (
		SELECT r.repo, i.import_path AS import
		FROM (SELECT id,CASE WHEN IFNULL(origin_url,'')='' THEN module_path ELSE origin_url END AS repo FROM repos) r
			 JOIN dirs d ON r.id=d.repo_id
			 JOIN files f ON d.id = f.dir_id
			 JOIN file_imports fi ON f.id = fi.file_id
			 JOIN imports i ON fi.import_id = i.id
		GROUP BY r.id,f.id,i.id 
	)x
	GROUP BY repo,import 
	ORDER BY repo,COUNT(import) DESC
`)
}

func displayRepoImports() {
	rows, err := queryRepoImports()
	if err != nil {
		log.Fatalf("Error querying database: %v\n", err)
	}
	defer mustClose(rows)
	repos := newImportGroups()
	repos.collect(rows)
	repos.display("repo")

}
