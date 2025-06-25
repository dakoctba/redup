package pkg

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Manager é responsável por gerenciar backups de arquivos duplicados
type Manager struct {
	backupDir string
	yes       bool
	logFile   string
}

// NewManager cria uma nova instância do gerenciador de backup
func NewManager(backupDir string, yes bool) *Manager {
	timestamp := time.Now().Format("20060102150405")
	logFile := fmt.Sprintf("redup-backup-%s.csv", timestamp)

	return &Manager{
		backupDir: backupDir,
		yes:       yes,
		logFile:   logFile,
	}
}

// ProcessDuplicates processa as duplicatas e move para backup
func (m *Manager) ProcessDuplicates(groups []FileGroup) error {
	if len(groups) == 0 {
		return nil
	}

	// Criar diretório de backup automaticamente
	backupPath, err := m.createBackupDirectory()
	if err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	fmt.Printf("Backup directory created: %s\n", backupPath)

	// Processar cada grupo de duplicatas
	for i, group := range groups {
		fmt.Printf("\nGroup %d:\n", i+1)

		// Mostrar lista numerada dos arquivos
		for j, file := range group.Files {
			fmt.Printf("[%d] %s\n", j+1, file.Path)
		}

		// Perguntar qual arquivo manter
		keepIndex := m.askWhichFileToKeep(len(group.Files))
		if keepIndex < 0 {
			fmt.Println("Skipping this group.")
			continue
		}

		// Mover todos os arquivos exceto o escolhido
		for j, file := range group.Files {
			if j == keepIndex {
				fmt.Printf("[%d] %s (keeping)\n", j+1, file.Path)
				continue
			}

			if m.confirmFileMove(file.Path) {
				// Passar o caminho do arquivo que será mantido
				keptFilePath := group.Files[keepIndex].Path
				if err := m.moveFileToBackup(file.Path, backupPath, keptFilePath, group.Checksum); err != nil {
					fmt.Printf("Error moving file %s: %v\n", file.Path, err)
				} else {
					fmt.Printf("→ Moved to %s\n", m.getBackupPath(file.Path, backupPath))
				}
			}
		}
	}

	return nil
}

// createBackupDirectory cria o diretório de backup
func (m *Manager) createBackupDirectory() (string, error) {
	timestamp := time.Now().Format("20060102150405")
	backupName := fmt.Sprintf("%s_backup", timestamp)
	backupPath := filepath.Join(m.backupDir, backupName)

	err := os.MkdirAll(backupPath, 0755)
	if err != nil {
		return "", err
	}

	return backupPath, nil
}

// confirmFileMove pergunta se deve mover um arquivo específico
func (m *Manager) confirmFileMove(filePath string) bool {
	// Se a flag --yes está ativada, move automaticamente
	if m.yes {
		return true
	}

	fmt.Printf("[y/N] Move duplicate: %s? ", filePath)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}

// moveFileToBackup move um arquivo para o diretório de backup
func (m *Manager) moveFileToBackup(filePath, backupPath, keptFilePath, checksum string) error {
	// Obter caminhos absolutos a partir da raiz do sistema
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", filePath, err)
	}

	absKeptFilePath, err := filepath.Abs(keptFilePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", keptFilePath, err)
	}

	// Garantir que os caminhos são absolutos a partir da raiz
	if !filepath.IsAbs(absFilePath) {
		absFilePath = "/" + absFilePath
	}
	if !filepath.IsAbs(absKeptFilePath) {
		absKeptFilePath = "/" + absKeptFilePath
	}

	// Criar a estrutura de diretórios no backup
	backupFilePath := m.getBackupPath(filePath, backupPath)
	backupDir := filepath.Dir(backupFilePath)

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory structure: %w", err)
	}

	// Mover o arquivo
	if err := os.Rename(filePath, backupFilePath); err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	// Obter caminho absoluto do backup também
	absBackupPath, err := filepath.Abs(backupFilePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for backup: %w", err)
	}
	if !filepath.IsAbs(absBackupPath) {
		absBackupPath = "/" + absBackupPath
	}

	// Adicionar entrada no arquivo CSV
	if err := m.addToCSV(absKeptFilePath, absFilePath, absBackupPath, checksum); err != nil {
		return fmt.Errorf("failed to add entry to CSV: %w", err)
	}

	return nil
}

// addToCSV adiciona uma entrada no arquivo CSV de backup
func (m *Manager) addToCSV(keptPath, movedPath, backupPath, checksum string) error {
	// Verificar se o arquivo CSV já existe
	fileExists := false
	if _, err := os.Stat(m.logFile); err == nil {
		fileExists = true
	}

	// Abrir arquivo para escrita (append se existir, criar se não existir)
	file, err := os.OpenFile(m.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Se o arquivo não existia, escrever cabeçalho
	if !fileExists {
		header := []string{"kept_path", "moved_path", "backup_path", "checksum", "timestamp"}
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	// Escrever linha de dados
	timestamp := time.Now().Format("2006-01-02T15:04:05")
	row := []string{keptPath, movedPath, backupPath, checksum, timestamp}

	return writer.Write(row)
}

// getBackupPath calcula o caminho do arquivo no diretório de backup
func (m *Manager) getBackupPath(filePath, backupPath string) string {
	// Se o arquivo for relativo, usar como está
	if !filepath.IsAbs(filePath) {
		return filepath.Join(backupPath, filePath)
	}

	// Se for absoluto, preservar a estrutura completa
	return filepath.Join(backupPath, filePath)
}

// askWhichFileToKeep pergunta ao usuário qual arquivo manter
func (m *Manager) askWhichFileToKeep(fileCount int) int {
	// Se a flag --yes está ativada, manter o primeiro arquivo automaticamente
	if m.yes {
		return 0
	}

	for {
		fmt.Printf("Which file to keep? (1-%d): ", fileCount)

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Permitir sair digitando 'q' ou 'quit'
		if input == "q" || input == "quit" {
			return -1
		}

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("Invalid input. Please enter a number between 1 and %d, or 'q' to quit.\n", fileCount)
			continue
		}

		if choice < 1 || choice > fileCount {
			fmt.Printf("Invalid choice. Please enter a number between 1 and %d.\n", fileCount)
			continue
		}

		// Retornar índice baseado em 0
		return choice - 1
	}
}
