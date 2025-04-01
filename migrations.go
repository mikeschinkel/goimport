package main

import (
	"embed"
	"fmt"
)

//go:embed migrations/*.sql
var migrationSQL embed.FS

// MustLoad reads a template file from the embed
//
//	Usage:
//		content, err := mustReadTemplate("specific.tmpl")
func mustLoadMigrationSQL(name string) (sql string) {
	b, err := migrationSQL.ReadFile(name)
	if err != nil {
		panic(fmt.Sprintf("%s not a valid migrations SQL; %s", name, err.Error()))
	}
	return string(b)
}
