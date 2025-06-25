package pkg

import (
	"fmt"
	"sort"
	"strings"
)

// FileGroup representa um grupo de arquivos com o mesmo checksum
type FileGroup struct {
	Checksum string
	Files    []FileInfo
	Size     int64
}

// DeduplicatorHasher é responsável por agrupar arquivos por checksum
type DeduplicatorHasher struct {
	hasher *Hasher
}

// NewDeduplicatorHasher cria uma nova instância do hasher para deduplicação
func NewDeduplicatorHasher(algorithm string) *DeduplicatorHasher {
	return &DeduplicatorHasher{
		hasher: NewHasher(algorithm),
	}
}

// GroupByChecksum agrupa arquivos por checksum
func (h *DeduplicatorHasher) GroupByChecksum(files []FileInfo) ([]FileGroup, error) {
	checksumMap := make(map[string][]FileInfo)

	for _, file := range files {
		// Log dinâmico na mesma linha - limpa completamente a linha anterior
		fmt.Printf("\r\033[KAnalisando: %s", file.Path)
		checksum, err := h.hasher.CalculateChecksum(file.Path)
		if err != nil {
			fmt.Println() // Garante nova linha em caso de erro
			return nil, fmt.Errorf("failed to calculate checksum for %s: %w", file.Path, err)
		}

		checksumMap[checksum] = append(checksumMap[checksum], file)
	}
	fmt.Println() // Nova linha ao terminar

	var groups []FileGroup
	for checksum, fileList := range checksumMap {
		if len(fileList) > 0 {
			// Ordenar arquivos por data de modificação e nome
			sort.Slice(fileList, func(i, j int) bool {
				// Se as datas são iguais, arquivos com "copy" são considerados mais novos
				if fileList[i].ModTime.Equal(fileList[j].ModTime) {
					// Arquivo com "copy" no nome vai para o final (será movido)
					iHasCopy := containsIgnoreCase(fileList[i].Path, "copy")
					jHasCopy := containsIgnoreCase(fileList[j].Path, "copy")

					// Se i tem "copy" e j não, i vai depois (mais novo)
					if iHasCopy && !jHasCopy {
						return false
					}
					// Se j tem "copy" e i não, j vai depois (mais novo)
					if !iHasCopy && jHasCopy {
						return true
					}
				}

				// Ordenação normal por data (mais antigo primeiro)
				return fileList[i].ModTime.Before(fileList[j].ModTime)
			})

			groups = append(groups, FileGroup{
				Checksum: checksum,
				Files:    fileList,
				Size:     fileList[0].Size, // Todos os arquivos no grupo têm o mesmo tamanho
			})
		}
	}

	return groups, nil
}

// FilterDuplicates retorna apenas grupos que contêm duplicatas (mais de um arquivo)
func FilterDuplicates(groups []FileGroup) []FileGroup {
	var duplicates []FileGroup

	for _, group := range groups {
		if len(group.Files) > 1 {
			duplicates = append(duplicates, group)
		}
	}

	return duplicates
}

// GetTotalDuplicateSize calcula o tamanho total que pode ser liberado removendo duplicatas
func GetTotalDuplicateSize(groups []FileGroup) int64 {
	var total int64

	for _, group := range groups {
		if len(group.Files) > 1 {
			// Calcular espaço que pode ser liberado (todos menos um arquivo)
			total += group.Size * int64(len(group.Files)-1)
		}
	}

	return total
}

// containsIgnoreCase verifica se uma string contém outra, ignorando maiúsculas/minúsculas
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
