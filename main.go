package main

import (
	_ "github.com/mattn/go-sqlite3"
)

const (
	DBSchemaVersion = 2
	configSubdir    = "goimports"
	defaultDBFile   = "goimports.db"
)

var inputs = initializeInputs()

func main() {
	ensureDB()
	defer mustClose(db)
	defer mustRollback(tx)
	ids := loadIdsFromInputs()
	updateDB(ids)
	displayOutput(ids)
	mustCommit(tx)
}
