package deduplicator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dakoctba/redup/scanner"
)

func TestNewHasher(t *testing.T) {
	hasher := NewHasher("sha256")
	if hasher.hasher.GetAlgorithm() != "sha256" {
		t.Errorf("Expected algorithm 'sha256', got '%s'", hasher.hasher.GetAlgorithm())
	}
}

func TestFilterDuplicates(t *testing.T) {
	// Criar arquivos temporários para teste
	tmpDir, err := os.MkdirTemp("", "test_dedup")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Criar arquivos com conteúdo idêntico (duplicatas)
	content1 := "duplicate content"
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	file3 := filepath.Join(tmpDir, "file3.txt")

	// Criar arquivo único
	content2 := "unique content"
	file4 := filepath.Join(tmpDir, "file4.txt")

	// Escrever conteúdo nos arquivos
	files := []string{file1, file2, file3, file4}
	contents := []string{content1, content1, content1, content2}

	for i, file := range files {
		err := os.WriteFile(file, []byte(contents[i]), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file %s: %v", file, err)
		}
	}

	// Criar FileInfo structs
	var fileInfos []scanner.FileInfo
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			t.Fatalf("Failed to stat file %s: %v", file, err)
		}
		fileInfos = append(fileInfos, scanner.FileInfo{
			Path: file,
			Size: info.Size(),
		})
	}

	// Agrupar por checksum
	hasher := NewHasher("sha256")
	groups, err := hasher.GroupByChecksum(fileInfos)
	if err != nil {
		t.Fatalf("Failed to group by checksum: %v", err)
	}

	// Filtrar duplicatas
	duplicates := FilterDuplicates(groups)

	// Verificar se encontrou apenas um grupo de duplicatas (3 arquivos idênticos)
	if len(duplicates) != 1 {
		t.Errorf("Expected 1 duplicate group, got %d", len(duplicates))
	}

	if len(duplicates[0].Files) != 3 {
		t.Errorf("Expected 3 duplicate files, got %d", len(duplicates[0].Files))
	}
}

func TestGetTotalDuplicateSize(t *testing.T) {
	// Criar grupos de teste
	groups := []FileGroup{
		{
			Checksum: "hash1",
			Files: []scanner.FileInfo{
				{Path: "file1.txt", Size: 100},
				{Path: "file2.txt", Size: 100},
				{Path: "file3.txt", Size: 100},
			},
			Size: 100,
		},
		{
			Checksum: "hash2",
			Files: []scanner.FileInfo{
				{Path: "file4.txt", Size: 200},
				{Path: "file5.txt", Size: 200},
			},
			Size: 200,
		},
		{
			Checksum: "hash3",
			Files: []scanner.FileInfo{
				{Path: "file6.txt", Size: 50},
			},
			Size: 50,
		},
	}

	// Calcular tamanho total que pode ser liberado
	totalSize := GetTotalDuplicateSize(groups)

	// Grupo 1: 2 arquivos podem ser removidos (200 bytes)
	// Grupo 2: 1 arquivo pode ser removido (200 bytes)
	// Grupo 3: 0 arquivos podem ser removidos (0 bytes)
	expectedSize := int64(400)

	if totalSize != expectedSize {
		t.Errorf("Expected total size %d, got %d", expectedSize, totalSize)
	}
}

func TestFilterDuplicatesEmpty(t *testing.T) {
	groups := []FileGroup{}
	duplicates := FilterDuplicates(groups)
	if len(duplicates) != 0 {
		t.Errorf("Expected 0 duplicate groups, got %d", len(duplicates))
	}
}

func TestFilterDuplicatesNoDuplicates(t *testing.T) {
	groups := []FileGroup{
		{
			Checksum: "hash1",
			Files: []scanner.FileInfo{
				{Path: "file1.txt", Size: 100},
			},
			Size: 100,
		},
		{
			Checksum: "hash2",
			Files: []scanner.FileInfo{
				{Path: "file2.txt", Size: 200},
			},
			Size: 200,
		},
	}

	duplicates := FilterDuplicates(groups)
	if len(duplicates) != 0 {
		t.Errorf("Expected 0 duplicate groups, got %d", len(duplicates))
	}
}
