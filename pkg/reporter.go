package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// PrintSummary exibe um resumo das duplicatas encontradas
func PrintSummary(groups []FileGroup) {
	if len(groups) == 0 {
		fmt.Println("No duplicate files found.")
		return
	}

	// Contar total de arquivos duplicados (excluindo o original de cada grupo)
	totalDuplicates := 0
	for _, group := range groups {
		totalDuplicates += len(group.Files) - 1
	}

	fmt.Printf("Found %d duplicate files:\n\n", totalDuplicates)

	// Mostrar cada grupo de duplicatas
	for i, group := range groups {
		if len(group.Files) > 1 {
			fmt.Printf("[%d] %s\n", i+1, group.Files[0].Path)
			fmt.Printf("Found %d copies:\n", len(group.Files)-1)

			// Mostrar todas as cópias (excluindo o primeiro arquivo que é considerado o original)
			for j := 1; j < len(group.Files); j++ {
				fmt.Printf("  %s\n", group.Files[j].Path)
			}
			fmt.Println()
		}
	}

	// Mostrar estatísticas
	totalSize := GetTotalDuplicateSize(groups)
	fmt.Printf("Total space that can be freed: %s\n", formatBytes(totalSize))
}

// ExportJSON exporta os resultados em formato JSON
func ExportJSON(groups []FileGroup, writer io.Writer) error {
	type FileInfo struct {
		Path string `json:"path"`
		Size int64  `json:"size"`
	}

	type GroupInfo struct {
		Checksum string     `json:"checksum"`
		Size     int64      `json:"size"`
		Files    []FileInfo `json:"files"`
	}

	var jsonGroups []GroupInfo
	for _, group := range groups {
		var files []FileInfo
		for _, file := range group.Files {
			files = append(files, FileInfo{
				Path: file.Path,
				Size: file.Size,
			})
		}

		jsonGroups = append(jsonGroups, GroupInfo{
			Checksum: group.Checksum,
			Size:     group.Size,
			Files:    files,
		})
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jsonGroups)
}

// ExportCSV exporta os resultados em formato CSV
func ExportCSV(groups []FileGroup, writer io.Writer) error {
	fmt.Fprintln(writer, "Group,Checksum,Size,File")

	for i, group := range groups {
		for _, file := range group.Files {
			fmt.Fprintf(writer, "%d,%s,%d,%s\n",
				i+1,
				group.Checksum,
				group.Size,
				strings.ReplaceAll(file.Path, ",", "\\,"))
		}
	}

	return nil
}

// formatBytes formata bytes em uma string legível
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
