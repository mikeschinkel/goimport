package main

import (
	"fmt"
	"log"
	"strings"
)

func displayOutput(ids Ids) {
	switch inputs.outputMode {
	case NoneOutput:
		// Display nothing
	case AllOutput:
		displayAllImports(ids)
	case DirsOutput:
		displayDirImports()
	case ReposOutput:
		displayRepoImports()
	case FileImportsOutput:
		displayFileImports(ids.dirId)
	case ImportFilesOutput:
		displayImportFiles(ids.dirId)
	default:
		log.Fatalf("Invalid mode: %s. Use 'files' or 'imports'\n", inputs.outputMode)
	}
}

func displayAllImports(ids Ids) {
	rows, err := queryAllImports(inputs.sort)
	if err != nil {
		log.Fatalf("Error querying database (dir_id=%d_: %v\n", ids.dirId, err)
	}
	defer mustClose(rows)
	fmt.Printf("All Imports Across %d Repos\n", countRepos())
	fmt.Println("Files — Import Path")
	fmt.Println("===================")
	for rows.Next() {
		var importPath string
		var fileCount int
		err = rows.Scan(&importPath, &fileCount)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		fmt.Printf("%5d — %s\n", fileCount, importPath)
	}
}

func displayFileImports(dirId int64) {
	rows, err := queryFileImportsForDirId(dirId)
	if err != nil {
		log.Fatalf("Error querying database (dir_id=%d): %v\n", dirId, err)
	}
	defer mustClose(rows)
	for rows.Next() {
		var filepath, imports string
		err = rows.Scan(&filepath, &imports)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		fmt.Printf("\n=== %s ===\n", filepath)
		importsList := strings.Split(imports, "\n")
		for _, imp := range importsList {
			fmt.Printf("import %s\n", imp)
		}
	}
}

func displayImportFiles(dirId int64) {
	rows, err := queryImportFilesForDirId(dirId)
	if err != nil {
		log.Fatalf("Error querying database (dir_id=%d): %v\n", dirId, err)
	}
	defer mustClose(rows)
	for rows.Next() {
		var importPath, filepaths string
		err = rows.Scan(&importPath, &filepaths)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		fmt.Printf("\n=== import %s ===\n", importPath)
		fmt.Printf("%s\n", filepaths)
	}
}
