package pkg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScannerIgnoresDotGit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_scanner_git")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Criar estrutura de arquivos
	os.Mkdir(filepath.Join(tmpDir, ".git"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("abc"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".git", "config"), []byte("should be ignored"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".git", "HEAD"), []byte("should be ignored"), 0644)

	scanner := NewScanner(0)
	files, err := scanner.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	for _, f := range files {
		if filepath.Base(f.Path) == "config" || filepath.Base(f.Path) == "HEAD" {
			t.Errorf("Arquivo do .git não deveria ser listado: %s", f.Path)
		}
	}

	// Deve encontrar apenas file1.txt
	if len(files) != 1 || filepath.Base(files[0].Path) != "file1.txt" {
		t.Errorf("Esperado apenas file1.txt, encontrado: %v", files)
	}
}
