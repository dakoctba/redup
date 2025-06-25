package pkg

// Config representa a configuração da aplicação
type Config struct {
	Dir       string
	Checksum  string
	MinSize   int64
	BackupDir string
	DryRun    bool
	JSON      bool
	Version   bool
	Yes       bool
}
