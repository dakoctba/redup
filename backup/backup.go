package backup

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dakoctba/redup/deduplicator"
)

// Manager é responsável por gerenciar backups de arquivos duplicados
type Manager struct {
	backupDir string
}

// NewManager cria uma nova instância do gerenciador de backup
func NewManager(backupDir string) *Manager {
	return &Manager{
		backupDir: backupDir,
	}
}

// ProcessDuplicates processa as duplicatas e move para backup
func (m *Manager) ProcessDuplicates(groups []deduplicator.FileGroup) error {
	if len(groups) == 0 {
		return nil
	}

	// Perguntar sobre criação do diretório de backup
	if !m.confirmBackupCreation() {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// Criar diretório de backup
	backupPath, err := m.createBackupDirectory()
	if err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	fmt.Printf("Backup directory created: %s\n", backupPath)

	// Processar cada grupo de duplicatas
	for i, group := range groups {
		fmt.Printf("\nGroup %d:\n", i+1)

		// Manter o primeiro arquivo, mover os demais
		for j, file := range group.Files {
			if j == 0 {
				fmt.Printf("[%d] %s (keeping)\n", j+1, file.Path)
				continue
			}

			if m.confirmFileMove(file.Path) {
				if err := m.moveFileToBackup(file.Path, backupPath); err != nil {
					fmt.Printf("Error moving file %s: %v\n", file.Path, err)
				} else {
					fmt.Printf("→ Moved to %s\n", m.getBackupPath(file.Path, backupPath))
				}
			}
		}
	}

	return nil
}

// confirmBackupCreation pergunta se deve criar o diretório de backup
func (m *Manager) confirmBackupCreation() bool {
	timestamp := time.Now().Format("20060102150405")
	backupName := fmt.Sprintf("%s_backup", timestamp)

	fmt.Printf("Create backup folder '%s'? [Y/n]: ", backupName)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "" || input == "y" || input == "yes"
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
	fmt.Printf("[y/N] Move duplicate: %s? ", filePath)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}

// moveFileToBackup move um arquivo para o diretório de backup
func (m *Manager) moveFileToBackup(filePath, backupPath string) error {
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

	return nil
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
