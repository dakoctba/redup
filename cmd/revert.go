package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	revertLogFile string
	revertDryRun  bool
)

// revertCmd represents the revert command
var revertCmd = &cobra.Command{
	Use:   "revert [log-file]",
	Short: "Revert files from backup using CSV log",
	Long: `Revert files from backup using the CSV log file generated during backup.
If no log file is specified, uses the most recent backup log.`,
	Example: `  redup revert                                    # Revert using most recent log
  redup revert redup-backup-20250625160855.csv  # Revert using specific log
  redup revert --dry-run                        # Simulate revert operation`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine log file to use
		logFile := revertLogFile
		if len(args) > 0 {
			logFile = args[0]
		}

		// If no log file specified, find the most recent one
		if logFile == "" {
			var err error
			logFile, err = findMostRecentLogFile()
			if err != nil {
				return fmt.Errorf("no backup log files found: %v", err)
			}
		}

		// Validate log file exists
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			return fmt.Errorf("log file '%s' does not exist", logFile)
		}

		fmt.Printf("Reverting from log file: %s\n", logFile)

		// Read and process CSV file
		if err := processRevert(logFile, revertDryRun); err != nil {
			return fmt.Errorf("error processing revert: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(revertCmd)

	revertCmd.Flags().StringVarP(&revertLogFile, "log", "l", "", "specific log file to use for revert")
	revertCmd.Flags().BoolVarP(&revertDryRun, "dry-run", "n", false, "simulate revert operation without moving files")
}

// findMostRecentLogFile finds the most recent backup log file
func findMostRecentLogFile() (string, error) {
	pattern := "redup-backup-*.csv"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no backup log files found")
	}

	// Find the most recent file (highest timestamp)
	var mostRecent string
	for _, match := range matches {
		if mostRecent == "" || match > mostRecent {
			mostRecent = match
		}
	}

	return mostRecent, nil
}

// processRevert processes the revert operation
func processRevert(logFile string, dryRun bool) error {
	// Open CSV file
	file, err := os.Open(logFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) < 2 {
		return fmt.Errorf("log file is empty or invalid")
	}

	// Skip header row
	records = records[1:]

	fmt.Printf("Found %d files to revert\n", len(records))

	successCount := 0
	errorCount := 0

	for i, record := range records {
		if len(record) < 3 {
			fmt.Printf("Warning: invalid record at line %d\n", i+2)
			errorCount++
			continue
		}

		movedPath := record[1]
		backupPath := record[2]

		if dryRun {
			fmt.Printf("[DRY-RUN] Would revert: %s -> %s\n", backupPath, movedPath)
			successCount++
			continue
		}

		// Check if backup file exists
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			fmt.Printf("Error: backup file not found: %s\n", backupPath)
			errorCount++
			continue
		}

		// Create directory structure for moved file
		movedDir := filepath.Dir(movedPath)
		if err := os.MkdirAll(movedDir, 0755); err != nil {
			fmt.Printf("Error creating directory for %s: %v\n", movedPath, err)
			errorCount++
			continue
		}

		// Move file from backup to original location
		if err := os.Rename(backupPath, movedPath); err != nil {
			fmt.Printf("Error moving file %s: %v\n", backupPath, err)
			errorCount++
			continue
		}

		fmt.Printf("Reverted: %s -> %s\n", backupPath, movedPath)
		successCount++
	}

	fmt.Printf("\nRevert completed: %d successful, %d errors\n", successCount, errorCount)

	// If not dry-run and all files were reverted successfully, remove the log file
	if !dryRun && errorCount == 0 && successCount > 0 {
		if err := os.Remove(logFile); err != nil {
			fmt.Printf("Warning: could not remove log file %s: %v\n", logFile, err)
		} else {
			fmt.Printf("Removed log file: %s\n", logFile)
		}

		// Try to remove backup directory if empty
		if len(records) > 0 {
			backupPath := records[0][2]
			backupDir := filepath.Dir(backupPath)
			if err := os.Remove(backupDir); err == nil {
				fmt.Printf("Removed empty backup directory: %s\n", backupDir)
			}
		}
	}

	return nil
}
