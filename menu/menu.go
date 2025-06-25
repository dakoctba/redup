package menu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dakoctba/redup/backup"
	"github.com/dakoctba/redup/deduplicator"
	"github.com/dakoctba/redup/reporter"
	"github.com/dakoctba/redup/scanner"
)

// Config representa a configuração da aplicação
type Config struct {
	Dir       string
	Checksum  string
	MinSize   int64
	BackupDir string
	DryRun    bool
	JSON      bool
}

// Menu representa o menu interativo
type Menu struct {
	config     *Config
	duplicates []deduplicator.FileGroup
	scanner    *scanner.Scanner
	hasher     *deduplicator.Hasher
	backupMgr  *backup.Manager
}

// NewMenu cria uma nova instância do menu
func NewMenu(config *Config) *Menu {
	return &Menu{
		config:    config,
		scanner:   scanner.NewScanner(config.MinSize),
		hasher:    deduplicator.NewHasher(config.Checksum),
		backupMgr: backup.NewManager(config.BackupDir),
	}
}

// Run executa o menu interativo
func (m *Menu) Run() {
	for {
		m.showMenu()
		choice := m.getChoice()

		switch choice {
		case 1:
			m.scanDirectory()
		case 2:
			m.showDuplicateSummary()
		case 3:
			m.removeDuplicates()
		case 4:
			m.exportResults()
		case 5:
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		fmt.Println()
	}
}

// showMenu exibe o menu principal
func (m *Menu) showMenu() {
	fmt.Println("redup — Duplicate File Manager")
	fmt.Println("[1] Scan directory")
	fmt.Println("[2] Show duplicate summary")
	fmt.Println("[3] Remove duplicates")
	fmt.Println("[4] Export results (JSON/CSV)")
	fmt.Println("[5] Exit")
}

// getChoice obtém a escolha do usuário
func (m *Menu) getChoice() int {
	fmt.Print("Choice: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	choice, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}

	return choice
}

// scanDirectory escaneia um diretório para duplicatas
func (m *Menu) scanDirectory() {
	if m.config.Dir == "." {
		fmt.Print("Enter directory path to scan: ")
		reader := bufio.NewReader(os.Stdin)
		dir, _ := reader.ReadString('\n')
		m.config.Dir = strings.TrimSpace(dir)
	}

	if m.config.Dir == "" {
		fmt.Println("No directory specified.")
		return
	}

	fmt.Printf("Scanning %s...\n", m.config.Dir)

	files, err := m.scanner.ScanDirectory(m.config.Dir)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return
	}

	fileGroups, err := m.hasher.GroupByChecksum(files)
	if err != nil {
		fmt.Printf("Error calculating checksums: %v\n", err)
		return
	}

	m.duplicates = deduplicator.FilterDuplicates(fileGroups)
	fmt.Printf("Found %d duplicate groups.\n", len(m.duplicates))
}

// showDuplicateSummary exibe o resumo das duplicatas
func (m *Menu) showDuplicateSummary() {
	if len(m.duplicates) == 0 {
		fmt.Println("No duplicates found. Please scan a directory first.")
		return
	}

	reporter.PrintSummary(m.duplicates)
}

// removeDuplicates remove as duplicatas
func (m *Menu) removeDuplicates() {
	if len(m.duplicates) == 0 {
		fmt.Println("No duplicates found. Please scan a directory first.")
		return
	}

	if err := m.backupMgr.ProcessDuplicates(m.duplicates); err != nil {
		fmt.Printf("Error processing duplicates: %v\n", err)
	}
}

// exportResults exporta os resultados
func (m *Menu) exportResults() {
	if len(m.duplicates) == 0 {
		fmt.Println("No duplicates found. Please scan a directory first.")
		return
	}

	fmt.Println("Export format:")
	fmt.Println("[1] JSON")
	fmt.Println("[2] CSV")

	choice := m.getChoice()

	switch choice {
	case 1:
		if err := reporter.ExportJSON(m.duplicates, os.Stdout); err != nil {
			fmt.Printf("Error exporting JSON: %v\n", err)
		}
	case 2:
		if err := reporter.ExportCSV(m.duplicates, os.Stdout); err != nil {
			fmt.Printf("Error exporting CSV: %v\n", err)
		}
	default:
		fmt.Println("Invalid choice.")
	}
}
