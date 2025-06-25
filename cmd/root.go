package cmd

import (
	"fmt"
	"os"

	"github.com/dakoctba/redup/pkg"
	"github.com/spf13/cobra"
)

var (
	// Versão padrão que será substituída pelo GoReleaser durante a compilação
	version   = "unknown"
	buildTime = "unknown"
	gitCommit = "unknown"

	// Flags
	dir       string
	checksum  string
	minSize   int64
	backupDir string
	dryRun    bool
	json      bool
	yes       bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "redup [directory]",
	Short: "Duplicate File Manager - Find and manage duplicate files by content",
	Long: `redup is a command line tool that allows you to find and manage
duplicate files by content using checksums (SHA-256 or MD5),
respecting .gitignore rules and providing safe backup options.`,
	Example: `  redup --dir ~/Documents --min-size 1048576    # Scan with minimum size
  redup --checksum md5 --dry-run ~/Pictures      # Use MD5, dry run
  redup --json ~/Music > duplicates.json         # Export to JSON
  redup --backup-dir ~/backups ~/Downloads       # Custom backup directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Se não há argumentos e nenhuma flag específica foi usada, mostrar ajuda
		if len(args) == 0 && dir == "." && checksum == "sha256" && minSize == 0 &&
			backupDir == "." && !dryRun && !json {
			return cmd.Help()
		}

		// Determine directory to scan
		scanDir := dir
		if len(args) > 0 {
			scanDir = args[0]
		}

		// Validate directory
		if _, err := os.Stat(scanDir); os.IsNotExist(err) {
			return fmt.Errorf("directory '%s' does not exist", scanDir)
		}

		// Configure processor
		config := pkg.Config{
			Dir:       scanDir,
			Checksum:  checksum,
			MinSize:   minSize,
			BackupDir: backupDir,
			DryRun:    dryRun,
			JSON:      json,
			Yes:       yes,
		}

		// Scan directory
		fmt.Printf("Scanning %s...\n", scanDir)

		fileScanner := pkg.NewScanner(config.MinSize)
		files, err := fileScanner.ScanDirectory(config.Dir)
		if err != nil {
			return fmt.Errorf("error scanning directory: %v", err)
		}

		// Calculate checksums
		hasher := pkg.NewDeduplicatorHasher(config.Checksum)
		fileGroups, err := hasher.GroupByChecksum(files)
		if err != nil {
			return fmt.Errorf("error calculating checksums: %v", err)
		}

		// Filter only duplicate groups
		duplicateGroups := pkg.FilterDuplicates(fileGroups)

		// Display results
		if config.JSON {
			pkg.ExportJSON(duplicateGroups, os.Stdout)
		} else {
			pkg.PrintSummary(duplicateGroups)
		}

		if len(duplicateGroups) == 0 {
			fmt.Println("No duplicate files found.")
			return nil
		}

		// If not dry-run, ask about backup
		if !config.DryRun {
			backupManager := pkg.NewManager(config.BackupDir, config.Yes)
			if err := backupManager.ProcessDuplicates(duplicateGroups); err != nil {
				return fmt.Errorf("error processing duplicates: %v", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&dir, "dir", "d", ".", "directory to scan (default: current working directory)")
	rootCmd.Flags().StringVarP(&checksum, "checksum", "c", "sha256", "checksum algorithm (sha256|md5)")
	rootCmd.Flags().Int64VarP(&minSize, "min-size", "s", 0, "minimum file size to consider in bytes")
	rootCmd.Flags().StringVarP(&backupDir, "backup-dir", "b", ".", "base directory for backup")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "simulate actions without moving files")
	rootCmd.Flags().BoolVarP(&json, "json", "j", false, "output results in JSON format")
	rootCmd.Flags().BoolVarP(&yes, "yes", "y", false, "move automatically all duplicates without asking for confirmation")

	rootCmd.Flags().BoolP("version", "v", false, "Show version number")

	rootCmd.SetVersionTemplate(`{{.Name}} version {{.Version}}
build time: ` + buildTime + `
git commit: ` + gitCommit + `
`)
	rootCmd.Version = version

	// Permitir que o comando completion apareça
	rootCmd.CompletionOptions.DisableDefaultCmd = false

	// Configurar um template de ajuda personalizado
	const customHelpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

	const customUsageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (and (not .Hidden) (ne .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

	// Aplicar os templates personalizados
	rootCmd.SetHelpTemplate(customHelpTemplate)
	rootCmd.SetUsageTemplate(customUsageTemplate)

	// Remover o comando help padrão
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "hidden-help",
		Hidden: true,
	})

	// Adicionar comando help oculto
	helpCmd := &cobra.Command{
		Use:    "help",
		Short:  "Help about any command",
		Hidden: true,
		Run: func(c *cobra.Command, args []string) {
			rootCmd.Help()
		},
	}
	rootCmd.AddCommand(helpCmd)

	// Adicionar o comando version, mas oculto
	versionCmd := &cobra.Command{
		Use:    "version",
		Short:  "Display application version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("redup version %s\n", version)
			fmt.Printf("build time: %s\n", buildTime)
			fmt.Printf("git commit: %s\n", gitCommit)
		},
	}
	rootCmd.AddCommand(versionCmd)
}

// SetVersionInfo permite que o main.go injete as informações de versão
func SetVersionInfo(v, bt, gc string) {
	version = v
	buildTime = bt
	gitCommit = gc
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
