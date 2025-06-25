package pkg

import (
	"os"
	"testing"
)

func TestNewHasher(t *testing.T) {
	hasher := NewHasher("sha256")
	if hasher.algorithm != "sha256" {
		t.Errorf("Expected algorithm 'sha256', got '%s'", hasher.algorithm)
	}
}

func TestCalculateChecksum(t *testing.T) {
	// Criar um arquivo temporário para teste
	content := "test content for checksum calculation"
	tmpFile, err := os.CreateTemp("", "test_checksum")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Testar SHA256
	hasher := NewHasher("sha256")
	checksum, err := hasher.CalculateChecksum(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to calculate SHA256 checksum: %v", err)
	}

	// Verificar se o checksum não está vazio
	if checksum == "" {
		t.Error("Expected non-empty checksum")
	}

	// Verificar se o checksum tem o tamanho correto para SHA256 (64 caracteres hex)
	if len(checksum) != 64 {
		t.Errorf("Expected SHA256 checksum length 64, got %d", len(checksum))
	}
}

func TestCalculateChecksumMD5(t *testing.T) {
	// Criar um arquivo temporário para teste
	content := "test content for MD5 checksum"
	tmpFile, err := os.CreateTemp("", "test_md5")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Testar MD5
	hasher := NewHasher("md5")
	checksum, err := hasher.CalculateChecksum(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to calculate MD5 checksum: %v", err)
	}

	// Verificar se o checksum não está vazio
	if checksum == "" {
		t.Error("Expected non-empty checksum")
	}

	// Verificar se o checksum tem o tamanho correto para MD5 (32 caracteres hex)
	if len(checksum) != 32 {
		t.Errorf("Expected MD5 checksum length 32, got %d", len(checksum))
	}
}

func TestCalculateChecksumUnsupportedAlgorithm(t *testing.T) {
	hasher := NewHasher("unsupported")
	_, err := hasher.CalculateChecksum("nonexistent")
	if err == nil {
		t.Error("Expected error for unsupported algorithm")
	}
}

func TestCalculateChecksumNonexistentFile(t *testing.T) {
	hasher := NewHasher("sha256")
	_, err := hasher.CalculateChecksum("nonexistent_file.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestGetAlgorithm(t *testing.T) {
	algorithm := "sha256"
	hasher := NewHasher(algorithm)
	if hasher.GetAlgorithm() != algorithm {
		t.Errorf("Expected algorithm '%s', got '%s'", algorithm, hasher.GetAlgorithm())
	}
}
