package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

var db *sql.DB
var tx *sql.Tx

func getDB() *sql.DB {
	return db
}
func getTx() *sql.Tx {
	return tx
}

// Check database version
func ensureDB() {
	var err error
	// Open or create the database
	db, err = sql.Open("sqlite3", inputs.dbPath)
	if err != nil {
		log.Fatalf("Error opening database %s: %v\n", inputs.dbPath, err)
	}
	tx, err = db.Begin()
	if err != nil {
		log.Fatalf("Error beginning transaction: %v\n", err)
	}
	ensureLatestDatabaseVersion()
}

func updateDB(ids Ids) {
	now := selectTimeNow()
	updateImportsForIds(ids)
	maybeInsertImports(ids)
	deleteFileImportsForDirIdOlderThan(ids.dirId, now)
	deleteFilesForDirIdOlderThan(ids.dirId, now)
}

// Check database version
func ensureLatestDatabaseVersion() {
	// Check DB version and migrate if needed
	dbVersion := getDatabaseVersion()
	if dbVersion == 0 {
		createInitialDatabase(db)
	}
	switch {
	case dbVersion < DBSchemaVersion:
		migrateDatabase(dbVersion)
	case dbVersion > DBSchemaVersion:
		log.Fatalf("Database schema version (%d) is greater than the version this app expects (%d); cannot continue", dbVersion, DBSchemaVersion)
	default:
		logf("Using database schema version %d", dbVersion)
	}
}

// Check if version table exists and what version it reports
func getDatabaseVersion() (version int) {
	var err error

	if !tableExists("version") && tableExists("imports") {
		// Must be version 1 since we don't have version but do have imports table
		version = 1
		goto end
	}
	err = db.QueryRow("SELECT version FROM version").Scan(&version)
	if err == nil {
		// Version was found and scanned
		goto end
	}
	if !tableExists("imports%") {
		// Neither imports or imports_old exist so version is 0 aka DB has never been
		// created
		goto end
	}
	// Assume version 1 if imports table exists
	version = 1
end:
	return version
}

// Migrate database from oldVersion to DBSchemaVersion
func migrateDatabase(oldVersion int) {
	var err error

	if oldVersion == DBSchemaVersion {
		goto end
	}

	logf("Migrating database from version %d to %d\n", oldVersion, DBSchemaVersion)

	for v := oldVersion + 1; v <= DBSchemaVersion; v++ {
		// Migrate from version 1 to 2
		query := mustLoadMigrationSQL(fmt.Sprintf(`migrations/version_%d.sql`, v))
		_, err = tx.Exec(query)
		if err != nil {
			log.Fatalf("Error migrating to schema version %d: %v\n", oldVersion, err)
		}
	}
	// Delete any existing version entries and insert the new version
	_, err = tx.Exec(`
		-- Delete any existing values to ensure we have only one version row
		DELETE FROM version;
		
		-- Update schema version
		INSERT INTO version (version) VALUES ($1);
	`, DBSchemaVersion)
	if err != nil {
		log.Fatalf("Error updating schema version: %v\n", err)
	}
	mustCommit(tx)
	log.Println("Database migration completed successfully")
	log.Println("Please re-run goimports now")
	os.Exit(0)
end:
}

// Check if version table exists and what version it reports
// Check if a table exists
func tableExists(name string) (exists bool) {
	var result int
	where := "= ?"
	if strings.Contains(name, "%") {
		where = "LIKE ?"
	}

	// Use EXISTS to properly check if the table exists
	query := fmt.Sprintf(`SELECT count(*) FROM sqlite_master WHERE type='table' AND name %s`, where)

	err := tx.QueryRow(query, name).Scan(&result)
	if err != nil {
		// TODO Inspect sentinel error here to handle fail on non-expected errors
		goto end
	}

	// result will be 1 if exists, 0 if not
	exists = result == 1
end:
	return exists
}

func selectTimeNow() (now string) {
	err := db.QueryRow("SELECT DATETIME('now');").Scan(&now)
	if err != nil {
		log.Fatal("Error selecting datetime: ", err)
	}
	return now
}
