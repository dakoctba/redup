# Makefile para o projeto Redup

# Variáveis
BINARY_NAME=redup
VERSION=$(shell git tag -l --sort=-v:refname | head -n 1 | sed 's/^v//' 2>/dev/null || echo "dev")
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/dakoctba/redup/cmd.version=$(VERSION) -X github.com/dakoctba/redup/cmd.buildTime=$(BUILD_TIME) -X github.com/dakoctba/redup/cmd.gitCommit=$(GIT_COMMIT)"

# Detecção do sistema operacional
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	OS := darwin
else ifeq ($(UNAME_S),Linux)
	OS := linux
else
	OS := windows
endif

# Detecção da arquitetura
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),x86_64)
	ARCH := amd64
else ifeq ($(UNAME_M),arm64)
	ARCH := arm64
else ifeq ($(UNAME_M),aarch64)
	ARCH := arm64
else
	ARCH := 386
endif

.PHONY: all build clean test run install uninstall release snapshot help

# Alvos padrão
all: build

# Compilação local para desenvolvimento
build:
	@echo "Building $(BINARY_NAME) $(VERSION) ($(GIT_COMMIT)) for $(OS)/$(ARCH)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME)
	@echo "Build completed successfully!"

# Limpar arquivos de build
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@go clean
	@echo "Clean completed successfully!"

# Executar testes
test:
	@echo "Running tests..."
	@go test ./...
	@echo "Tests completed successfully!"

# Executar o programa (para desenvolvimento)
run:
	@go run $(LDFLAGS) main.go $(ARGS)

# Instalação local
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS)
	@echo "Installation completed successfully!"

# Desinstalação
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(shell which $(BINARY_NAME))
	@echo "Uninstallation completed successfully!"

# Processo interativo para criar uma nova versão e release
release:
	@echo "Iniciando processo interativo de criação de nova versão e release..."
	@bin/update_version.sh

# Criar uma release em modo snapshot (para testes)
snapshot:
	@echo "Creating snapshot release..."
	@bin/release.sh --snapshot

# Exibir informações de ajuda
help:
	@echo "Makefile para $(BINARY_NAME) - Comandos disponíveis:"
	@echo ""
	@echo "  make                - Equivalente a 'make build'"
	@echo "  make build          - Compilar o projeto para o ambiente local"
	@echo "  make clean          - Remover arquivos temporários e de build"
	@echo "  make test           - Executar testes"
	@echo "  make run ARGS=\"\"   - Executar aplicação (passar argumentos em ARGS)"
	@echo "  make install        - Instalar o binário localmente"
	@echo "  make uninstall      - Desinstalar o binário"
	@echo "  make release        - Processo interativo para criar nova versão e release"
	@echo "  make snapshot       - Criar uma release de teste (não publicada)"
	@echo "  make help           - Exibir esta mensagem de ajuda"
	@echo ""
	@echo "Exemplos:"
	@echo "  make run ARGS=\"--dry-run .\""
	@echo "  make run ARGS=\"--version\""
	@echo ""
