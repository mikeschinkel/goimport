package main

import (
	"database/sql"
	"fmt"
	"log"
)

type importGroup struct {
	path      string
	importMap map[string]*importInfo
	imports   []*importInfo
}

func newImportGroup(path string) *importGroup {
	return &importGroup{
		path:      path,
		importMap: make(map[string]*importInfo),
	}
}

type importGroups struct {
	groups   []*importGroup
	groupMap map[string]*importGroup
}

func (gg *importGroups) display(name string) {
	for _, d := range gg.groups {
		fmt.Printf("\nImports for %s %s:\n", name, d.path)
		imports := importInfos(d.imports)
		imports.sort(inputs.sort)
		for _, i := range imports {
			fmt.Printf("%5d — %s\n", i.count, i.path)
		}
	}
}

func (gg *importGroups) collect(rows *sql.Rows) {
	var err error
	for rows.Next() {
		var group, imp string
		var count int
		err = rows.Scan(&group, &imp, &count)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		gg.record(group, imp, count)
	}
}

func newImportGroups() *importGroups {
	return &importGroups{
		groupMap: make(map[string]*importGroup),
	}
}

func (dd *importGroups) record(dir, imp string, count int) {
	var d *importGroup
	d, ok := dd.groupMap[dir]
	if !ok {
		d = newImportGroup(dir)
		dd.groupMap[dir] = d
		dd.groups = append(dd.groups, d)
	}
	i, ok := d.importMap[imp]
	if !ok {
		i = &importInfo{path: imp, count: count}
		d.importMap[imp] = i
		d.imports = append(d.imports, i)
	}
}
