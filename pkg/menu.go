package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Menu representa o menu interativo
type Menu struct {
	config     *Config
	duplicates []FileGroup
	scanner    *Scanner
	hasher     *DeduplicatorHasher
	backupMgr  *Manager
}

// NewMenu cria uma nova instância do menu
func NewMenu(config *Config) *Menu {
	return &Menu{
		config:    config,
		scanner:   NewScanner(config.MinSize),
		hasher:    NewDeduplicatorHasher(config.Checksum),
		backupMgr: NewManager(config.BackupDir),
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
			m.showVersion()
		case 6:
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
	fmt.Println("[5] Show version")
	fmt.Println("[6] Exit")
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

	m.duplicates = FilterDuplicates(fileGroups)
	fmt.Printf("Found %d duplicate groups.\n", len(m.duplicates))
}

// showDuplicateSummary exibe o resumo das duplicatas
func (m *Menu) showDuplicateSummary() {
	if len(m.duplicates) == 0 {
		fmt.Println("No duplicates found. Please scan a directory first.")
		return
	}

	PrintSummary(m.duplicates)
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
		if err := ExportJSON(m.duplicates, os.Stdout); err != nil {
			fmt.Printf("Error exporting JSON: %v\n", err)
		}
	case 2:
		if err := ExportCSV(m.duplicates, os.Stdout); err != nil {
			fmt.Printf("Error exporting CSV: %v\n", err)
		}
	default:
		fmt.Println("Invalid choice.")
	}
}

// showVersion exibe informações da versão
func (m *Menu) showVersion() {
	fmt.Println("redup version dev")
	fmt.Println("Build time: unknown")
	fmt.Println("Git commit: unknown")
}
