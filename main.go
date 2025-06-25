package main

import (
	"log"
	"os"

	"github.com/dakoctba/redup/cmd"
	"github.com/joho/godotenv"
)

// Estas variáveis serão substituídas pelo Makefile durante a compilação
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if the .env file doesn't exist in production/deployment
		if _, err := os.Stat(".env"); !os.IsNotExist(err) {
			log.Printf("Warning: could not load .env file: %v", err)
		}
	}
}

func main() {
	// Passar as variáveis de versão para o cmd
	cmd.SetVersionInfo(Version, BuildTime, GitCommit)
	cmd.Execute()
}
