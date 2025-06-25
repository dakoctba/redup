package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GitignoreManager gerencia as regras do .gitignore
type GitignoreManager struct {
	rules []string
}

// NewGitignoreManager cria uma nova instância do gerenciador de .gitignore
func NewGitignoreManager() *GitignoreManager {
	return &GitignoreManager{
		rules: make([]string, 0),
	}
}

// LoadGitignore carrega as regras do arquivo .gitignore
func (g *GitignoreManager) LoadGitignore(rootDir string) error {
	gitignorePath := filepath.Join(rootDir, ".gitignore")

	file, err := os.Open(gitignorePath)
	if err != nil {
		// Se não existir .gitignore, não é um erro
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignorar linhas vazias e comentários
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		g.rules = append(g.rules, line)
	}

	return scanner.Err()
}

// ShouldIgnore verifica se um caminho deve ser ignorado
func (g *GitignoreManager) ShouldIgnore(path, rootDir string) bool {
	// Converter para caminho relativo ao diretório raiz
	relPath, err := filepath.Rel(rootDir, path)
	if err != nil {
		// Se não conseguir converter, não ignorar
		return false
	}

	// Normalizar separadores de caminho
	relPath = filepath.ToSlash(relPath)

	for _, rule := range g.rules {
		if g.matchesRule(relPath, rule) {
			return true
		}
	}

	return false
}

// matchesRule verifica se um caminho corresponde a uma regra do .gitignore
func (g *GitignoreManager) matchesRule(path, rule string) bool {
	rule = filepath.ToSlash(rule)

	// Se a regra começa e termina com /, é um diretório absoluto na raiz
	if strings.HasPrefix(rule, "/") && strings.HasSuffix(rule, "/") {
		dir := strings.Trim(rule, "/")
		return strings.HasPrefix(path, dir+"/") || path == dir
	}

	// Se a regra termina com /, é um diretório
	if strings.HasSuffix(rule, "/") {
		rule = rule[:len(rule)-1]
		return strings.HasPrefix(path, rule+"/") || path == rule
	}

	// Se a regra começa com /, deve corresponder ao início do caminho
	if strings.HasPrefix(rule, "/") {
		prefix := strings.TrimPrefix(rule, "/")
		return strings.HasPrefix(path, prefix)
	}

	// Se a regra contém *, é um padrão wildcard
	if strings.Contains(rule, "*") {
		return g.matchesWildcard(path, rule)
	}

	return path == rule || strings.HasSuffix(path, "/"+rule)
}

// matchesWildcard verifica se um caminho corresponde a um padrão wildcard
func (g *GitignoreManager) matchesWildcard(path, pattern string) bool {
	// Converter padrão simples para regex-like matching
	// Por exemplo: *.log -> qualquer arquivo que termine com .log

	// Se o padrão é *.ext
	if strings.HasPrefix(pattern, "*.") {
		ext := pattern[1:] // Remove o *
		return strings.HasSuffix(path, ext)
	}

	// Se o padrão termina com *
	if strings.HasSuffix(pattern, "*") {
		prefix := pattern[:len(pattern)-1]
		return strings.HasPrefix(path, prefix)
	}

	// Se o padrão contém * no meio
	if strings.Contains(pattern, "*") {
		parts := strings.Split(pattern, "*")
		if len(parts) == 2 {
			prefix := parts[0]
			suffix := parts[1]
			return strings.HasPrefix(path, prefix) && strings.HasSuffix(path, suffix)
		}
	}

	return false
}

// GetRules retorna as regras carregadas (para debug)
func (g *GitignoreManager) GetRules() []string {
	return g.rules
}
