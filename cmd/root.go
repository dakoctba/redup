package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/dakoctba/redup/pkg"
)

// Execute é a função principal que executa o comando
func Execute() int {
	config := parseFlags()

	if len(os.Args) == 1 {
		// Modo interativo
		return runInteractiveMode(config)
	} else {
		// Modo CLI
		return runCLIMode(config)
	}
}

func parseFlags() *pkg.Config {
	config := &pkg.Config{}

	flag.StringVar(&config.Dir, "dir", ".", "directory to scan (default: current working directory)")
	flag.StringVar(&config.Checksum, "checksum", "sha256", "checksum algorithm (sha256|md5)")
	flag.Int64Var(&config.MinSize, "min-size", 0, "minimum file size to consider in bytes")
	flag.StringVar(&config.BackupDir, "backup-dir", ".", "base directory for backup")
	flag.BoolVar(&config.DryRun, "dry-run", false, "simulate actions without moving files")
	flag.BoolVar(&config.JSON, "json", false, "output results in JSON format")

	flag.Parse()

	return config
}

func runInteractiveMode(config *pkg.Config) int {
	app := pkg.NewMenu(config)
	app.Run()
	return 0
}

func runCLIMode(config *pkg.Config) int {
	// Validar diretório
	if _, err := os.Stat(config.Dir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: directory '%s' does not exist\n", config.Dir)
		return 2
	}

	// Escanear diretório
	fmt.Printf("Scanning %s...\n", config.Dir)

	fileScanner := pkg.NewScanner(config.MinSize)
	files, err := fileScanner.ScanDirectory(config.Dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		return 2
	}

	// Calcular checksums
	hasher := pkg.NewDeduplicatorHasher(config.Checksum)
	fileGroups, err := hasher.GroupByChecksum(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calculating checksums: %v\n", err)
		return 2
	}

	// Filtrar apenas grupos com duplicatas
	duplicateGroups := pkg.FilterDuplicates(fileGroups)

	// Exibir resultados
	if config.JSON {
		pkg.ExportJSON(duplicateGroups, os.Stdout)
	} else {
		pkg.PrintSummary(duplicateGroups)
	}

	if len(duplicateGroups) == 0 {
		fmt.Println("No duplicate files found.")
		return 0
	}

	// Se não for dry-run, perguntar sobre backup
	if !config.DryRun {
		backupManager := pkg.NewManager(config.BackupDir)
		if err := backupManager.ProcessDuplicates(duplicateGroups); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing duplicates: %v\n", err)
			return 2
		}
	}

	return 0
}
