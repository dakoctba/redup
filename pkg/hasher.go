package pkg

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Hasher é responsável por calcular checksums de arquivos
type Hasher struct {
	algorithm string
}

// NewHasher cria uma nova instância do hasher
func NewHasher(algorithm string) *Hasher {
	return &Hasher{
		algorithm: algorithm,
	}
}

// CalculateChecksum calcula o checksum de um arquivo
func (h *Hasher) CalculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	var hash interface{}

	switch h.algorithm {
	case "sha256":
		hash = sha256.New()
	case "md5":
		hash = md5.New()
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", h.algorithm)
	}

	_, err = io.Copy(hash.(io.Writer), file)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return hex.EncodeToString(hash.(io.Writer).(interface{ Sum([]byte) []byte }).Sum(nil)), nil
}

// GetAlgorithm retorna o algoritmo atual
func (h *Hasher) GetAlgorithm() string {
	return h.algorithm
}
