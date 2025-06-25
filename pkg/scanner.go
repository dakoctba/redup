package pkg

import (
	"os"
	"path/filepath"
)

// FileInfo representa informações básicas de um arquivo
type FileInfo struct {
	Path string
	Size int64
}

// Scanner é responsável por escanear diretórios e encontrar arquivos
type Scanner struct {
	minSize      int64
	gitignoreMgr *GitignoreManager
	rootDir      string
}

// NewScanner cria uma nova instância do scanner
func NewScanner(minSize int64) *Scanner {
	return &Scanner{
		minSize:      minSize,
		gitignoreMgr: NewGitignoreManager(),
	}
}

// ScanDirectory escaneia recursivamente um diretório e retorna informações dos arquivos
func (s *Scanner) ScanDirectory(root string) ([]FileInfo, error) {
	s.rootDir = root

	// Carregar regras do .gitignore
	if err := s.gitignoreMgr.LoadGitignore(root); err != nil {
		return nil, err
	}

	var files []FileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorar o diretório .git
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Pular diretórios
		if info.IsDir() {
			return nil
		}

		// Verificar se o arquivo deve ser ignorado pelo .gitignore
		if s.gitignoreMgr.ShouldIgnore(path, s.rootDir) {
			return nil
		}

		// Verificar tamanho mínimo
		if s.minSize > 0 && info.Size() < s.minSize {
			return nil
		}

		// Adicionar arquivo à lista
		files = append(files, FileInfo{
			Path: path,
			Size: info.Size(),
		})

		return nil
	})

	return files, err
}

// ScanFromStdin lê uma lista de caminhos de arquivos do stdin
func (s *Scanner) ScanFromStdin() ([]FileInfo, error) {
	// Esta funcionalidade será implementada se necessário
	// Por enquanto, retorna erro indicando que não está implementada
	return nil, nil
}

// GetIgnoredRules retorna as regras do .gitignore carregadas (para debug)
func (s *Scanner) GetIgnoredRules() []string {
	return s.gitignoreMgr.GetRules()
}
