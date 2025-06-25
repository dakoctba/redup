package pkg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewGitignoreManager(t *testing.T) {
	manager := NewGitignoreManager()
	if manager == nil {
		t.Error("Expected non-nil GitignoreManager")
	}
	if len(manager.rules) != 0 {
		t.Error("Expected empty rules slice")
	}
}

func TestLoadGitignore(t *testing.T) {
	// Criar diretório temporário
	tmpDir, err := os.MkdirTemp("", "test_gitignore")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Criar arquivo .gitignore
	gitignoreContent := `# Comentário
*.log
*.tmp
node_modules/
.DS_Store
build/
*.exe
`
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	err = os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write .gitignore: %v", err)
	}

	// Testar carregamento
	manager := NewGitignoreManager()
	err = manager.LoadGitignore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load .gitignore: %v", err)
	}

	// Verificar se as regras foram carregadas (excluindo comentários e linhas vazias)
	expectedRules := []string{"*.log", "*.tmp", "node_modules/", ".DS_Store", "build/", "*.exe"}
	if len(manager.rules) != len(expectedRules) {
		t.Errorf("Expected %d rules, got %d", len(expectedRules), len(manager.rules))
	}
}

func TestLoadGitignoreNonexistent(t *testing.T) {
	// Criar diretório temporário sem .gitignore
	tmpDir, err := os.MkdirTemp("", "test_no_gitignore")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewGitignoreManager()
	err = manager.LoadGitignore(tmpDir)
	if err != nil {
		t.Errorf("Expected no error when .gitignore doesn't exist, got: %v", err)
	}

	if len(manager.rules) != 0 {
		t.Error("Expected empty rules when .gitignore doesn't exist")
	}
}

func TestShouldIgnore(t *testing.T) {
	// Criar diretório temporário
	tmpDir, err := os.MkdirTemp("", "test_ignore")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Criar arquivo .gitignore
	gitignoreContent := `*.log
*.tmp
node_modules/
.DS_Store
build/
*.exe
`
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	err = os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write .gitignore: %v", err)
	}

	manager := NewGitignoreManager()
	err = manager.LoadGitignore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load .gitignore: %v", err)
	}

	// Testar casos que devem ser ignorados
	shouldIgnoreCases := []string{
		filepath.Join(tmpDir, "app.log"),
		filepath.Join(tmpDir, "temp.tmp"),
		filepath.Join(tmpDir, "node_modules", "package.json"),
		filepath.Join(tmpDir, ".DS_Store"),
		filepath.Join(tmpDir, "build", "app.exe"),
	}

	for _, path := range shouldIgnoreCases {
		if !manager.ShouldIgnore(path, tmpDir) {
			t.Errorf("Expected %s to be ignored", path)
		}
	}

	// Testar casos que NÃO devem ser ignorados
	shouldNotIgnoreCases := []string{
		filepath.Join(tmpDir, "app.txt"),
		filepath.Join(tmpDir, "src", "main.go"),
		filepath.Join(tmpDir, "README.md"),
	}

	for _, path := range shouldNotIgnoreCases {
		if manager.ShouldIgnore(path, tmpDir) {
			t.Errorf("Expected %s NOT to be ignored", path)
		}
	}
}

func TestMatchesRule(t *testing.T) {
	manager := NewGitignoreManager()

	// Testar regras de arquivo
	if !manager.matchesRule("app.log", "*.log") {
		t.Error("Expected app.log to match *.log")
	}

	if manager.matchesRule("app.txt", "*.log") {
		t.Error("Expected app.txt NOT to match *.log")
	}

	// Testar regras de diretório
	if !manager.matchesRule("node_modules/package.json", "node_modules/") {
		t.Error("Expected node_modules/package.json to match node_modules/")
	}

	if manager.matchesRule("src/main.go", "node_modules/") {
		t.Error("Expected src/main.go NOT to match node_modules/")
	}

	// Testar regras com caminho absoluto
	if !manager.matchesRule("src/build/app.exe", "/src/build/") {
		t.Error("Expected src/build/app.exe to match /src/build/")
	}
}
