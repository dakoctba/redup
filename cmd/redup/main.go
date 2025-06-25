package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dakoctba/redup/backup"
	"github.com/dakoctba/redup/deduplicator"
	"github.com/dakoctba/redup/menu"
	"github.com/dakoctba/redup/reporter"
	"github.com/dakoctba/redup/scanner"
)

func main() {
	config := parseFlags()

	if len(os.Args) == 1 {
		// Modo interativo
		runInteractiveMode(config)
	} else {
		// Modo CLI
		runCLIMode(config)
	}
}

func parseFlags() *menu.Config {
	config := &menu.Config{}

	flag.StringVar(&config.Dir, "dir", ".", "directory to scan (default: current working directory)")
	flag.StringVar(&config.Checksum, "checksum", "sha256", "checksum algorithm (sha256|md5)")
	flag.Int64Var(&config.MinSize, "min-size", 0, "minimum file size to consider in bytes")
	flag.StringVar(&config.BackupDir, "backup-dir", ".", "base directory for backup")
	flag.BoolVar(&config.DryRun, "dry-run", false, "simulate actions without moving files")
	flag.BoolVar(&config.JSON, "json", false, "output results in JSON format")

	flag.Parse()

	return config
}

func runInteractiveMode(config *menu.Config) {
	app := menu.NewMenu(config)
	app.Run()
}

func runCLIMode(config *menu.Config) {
	// Validar diretório
	if _, err := os.Stat(config.Dir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: directory '%s' does not exist\n", config.Dir)
		os.Exit(2)
	}

	// Escanear diretório
	fmt.Printf("Scanning %s...\n", config.Dir)

	fileScanner := scanner.NewScanner(config.MinSize)
	files, err := fileScanner.ScanDirectory(config.Dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(2)
	}

	// Calcular checksums
	hasher := deduplicator.NewHasher(config.Checksum)
	fileGroups, err := hasher.GroupByChecksum(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calculating checksums: %v\n", err)
		os.Exit(2)
	}

	// Filtrar apenas grupos com duplicatas
	duplicateGroups := deduplicator.FilterDuplicates(fileGroups)

	// Exibir resultados
	if config.JSON {
		reporter.ExportJSON(duplicateGroups, os.Stdout)
	} else {
		reporter.PrintSummary(duplicateGroups)
	}

	if len(duplicateGroups) == 0 {
		fmt.Println("No duplicate files found.")
		os.Exit(0)
	}

	// Se não for dry-run, perguntar sobre backup
	if !config.DryRun {
		backupManager := backup.NewManager(config.BackupDir)
		if err := backupManager.ProcessDuplicates(duplicateGroups); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing duplicates: %v\n", err)
			os.Exit(2)
		}
	}
}
