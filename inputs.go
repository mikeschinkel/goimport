package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Inputs struct {
	rootDirFlag    *string
	outputModeFlag *string
	dbPathFlag     *string
	sortFlag       *string
	verboseFlag    *bool
	noUpdateFlag   *bool
	workingDir     string
	rootPath       string
	gitOriginURL   string
	modulePath     string
	dbPath         string
	outputMode     OutputMode
	sort           SortType
	verbose        bool
	noUpdate       bool
	repoId         int
}

type OutputMode string

const (
	AllOutput         OutputMode = "all"
	DirsOutput        OutputMode = "dira"
	ReposOutput       OutputMode = "repos"
	FileImportsOutput OutputMode = "files"
	ImportFilesOutput OutputMode = "imports"
	NoneOutput        OutputMode = "none"
)

type SortType string

const (
	FileSort   SortType = "files"
	ImportSort SortType = "imports"
	CountSort  SortType = "counts"
)

func getRootPath(inputs Inputs) string {
	var err error
	inputs.rootPath, err = findGitRoot(*inputs.rootDirFlag)
	if err != nil {
		// If not in a git repo, use the specified path
		inputs.rootPath = inputs.workingDir
	}
	inputs.rootPath, err = filepath.Abs(inputs.rootPath)
	if err != nil {
		log.Fatalf("Error getting absolute path of '%s': %v\n", inputs.rootPath, err)
	}
	return inputs.rootPath
}

func initializeInputs() (inputs Inputs) {

	log.SetOutput(os.Stderr)

	inputs = Inputs{
		noUpdateFlag:   flag.Bool("noupdate", false, "Bypassing updating the database"),
		verboseFlag:    flag.Bool("verbose", false, "Display verbose output (forced to true when -mode=none)"),
		rootDirFlag:    flag.String("dir", ".", "Root directory to scan"),
		outputModeFlag: flag.String("mode", "", "Display mode:\n\tfiles   — to show files with their imports,\n\timports — to show imports with their files,\n\tdirs    — to show imports by directory,\n\trepos   — to show imports by repo, or\n\tnone    — to show only database update/log messages"),
		dbPathFlag:     flag.String("db", "", "Custom SQLite database file path (default: ~/.inputs/goimports/goimports.db)"),
		sortFlag:       flag.String("sort", "files", "Sort for applicable display modes; can be 'files' or 'imports'"),
	}
	flag.Parse()

	inputs.noUpdate = *inputs.noUpdateFlag
	inputs.workingDir = getWorkingDir(*inputs.rootDirFlag)
	inputs.rootPath = getRootPath(inputs)
	inputs.gitOriginURL = getGitOrigin(inputs.rootPath)
	inputs.modulePath = getGoModule(inputs.rootPath)
	inputs.dbPath = getDBPath(*inputs.dbPathFlag)
	if *inputs.outputModeFlag == "" {
		fmt.Printf("goimports\n\n\tERROR: You must specify a valid -mode option.\n\nUsage:\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	inputs.outputMode = ensureValidOutput(*inputs.outputModeFlag)
	inputs.sort = ensureValidSort(*inputs.sortFlag, inputs.outputMode)
	if inputs.outputMode == NoneOutput {
		inputs.verbose = true
	}
	return inputs
}

func ensureValidOutput(outputMode string) OutputMode {
	switch OutputMode(outputMode) {
	case AllOutput, ReposOutput, DirsOutput, FileImportsOutput, ImportFilesOutput, NoneOutput:
		// S'all good
		goto end
	}
	log.Fatalf("Invalid output mode: %s", outputMode)
end:
	return OutputMode(outputMode)
}

func ensureValidSort(sort string, outputMode OutputMode) SortType {
	switch SortType(sort) {
	case FileSort, ImportSort, CountSort:
		// S'all good
		goto end
	}
	log.Fatalf("Invalid dir mode sort: %s", sort)
end:
	return SortType(sort)
}
