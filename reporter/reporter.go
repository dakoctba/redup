package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/dakoctba/redup/deduplicator"
)

// PrintSummary exibe um resumo das duplicatas encontradas
func PrintSummary(groups []deduplicator.FileGroup) {
	if len(groups) == 0 {
		fmt.Println("No duplicate files found.")
		return
	}

	fmt.Printf("Found %d duplicate groups.\n\n", len(groups))

	totalSize := deduplicator.GetTotalDuplicateSize(groups)
	fmt.Printf("Total space that can be freed: %s\n\n", formatBytes(totalSize))

	for i, group := range groups {
		fmt.Printf("Group %d (%s):\n", i+1, formatBytes(group.Size))
		for j, file := range group.Files {
			fmt.Printf("  [%d] %s\n", j+1, file.Path)
		}
		fmt.Println()
	}
}

// ExportJSON exporta os resultados em formato JSON
func ExportJSON(groups []deduplicator.FileGroup, writer io.Writer) error {
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
func ExportCSV(groups []deduplicator.FileGroup, writer io.Writer) error {
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
