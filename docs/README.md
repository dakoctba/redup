# redup - Duplicate File Manager

`redup` é um utilitário CLI em Go para detectar e gerenciar arquivos duplicados. Inspirado no projeto [scopy](https://github.com/dakoctba/scopy), o `redup` oferece uma interface interativa e de linha de comando para identificar arquivos duplicados por conteúdo (usando checksums SHA-256 ou MD5) e movê-los para um diretório de backup com estrutura preservada.

## Características

- **Detecção por Conteúdo**: Usa checksums SHA-256 ou MD5 para identificar duplicatas independentemente do nome do arquivo
- **Interface Dupla**: Modo CLI para automação e modo interativo para uso manual
- **Backup Seguro**: Move duplicatas para diretório de backup com timestamp, preservando a estrutura original
- **Filtros Configuráveis**: Tamanho mínimo de arquivo, algoritmo de checksum personalizável
- **Exportação**: Resultados em JSON ou CSV
- **Código Modular**: Arquitetura limpa e testável seguindo padrões do projeto scopy

## Instalação

### Via Go Install
```bash
go install github.com/dakoctba/redup/cmd/redup@latest
```

### Compilação Local
```bash
git clone https://github.com/dakoctba/redup.git
cd redup
go build -o redup cmd/redup/main.go
```

## Uso

### Modo Interativo
Execute sem argumentos para iniciar o menu interativo:
```bash
redup
```

O menu oferece as seguintes opções:
- `[1] Scan directory` - Escanear diretório para duplicatas
- `[2] Show duplicate summary` - Exibir resumo das duplicatas encontradas
- `[3] Remove duplicates` - Mover duplicatas para backup
- `[4] Export results` - Exportar resultados em JSON ou CSV
- `[5] Exit` - Sair do programa

### Modo CLI
```bash
# Escanear diretório atual
redup --dir .

# Escanear diretório específico com tamanho mínimo
redup --dir ~/Documents --min-size 1048576

# Usar MD5 em vez de SHA-256
redup --dir ~/Pictures --checksum md5

# Modo dry-run (apenas simular)
redup --dir ~/Downloads --dry-run

# Exportar resultados em JSON
redup --dir ~/Music --json

# Especificar diretório de backup
redup --dir ~/Videos --backup-dir ~/backups
```

### Flags Disponíveis

| Flag | Descrição | Padrão |
|------|-----------|--------|
| `--dir` | Diretório para escanear | `.` (diretório atual) |
| `--checksum` | Algoritmo de checksum (`sha256` ou `md5`) | `sha256` |
| `--min-size` | Tamanho mínimo de arquivo em bytes | `0` (sem limite) |
| `--backup-dir` | Diretório base para backup | `.` (diretório atual) |
| `--dry-run` | Simular ações sem mover arquivos | `false` |
| `--json` | Exportar resultados em JSON | `false` |
| `--help` | Mostrar ajuda | - |

## Como Funciona

### 1. Detecção de Duplicatas
O `redup` escaneia recursivamente o diretório especificado e:
- Filtra arquivos por tamanho mínimo (se especificado)
- Calcula checksums usando o algoritmo escolhido
- Agrupa arquivos com checksums idênticos

### 2. Processamento de Duplicatas
Quando o usuário escolhe remover duplicatas:
1. Solicita confirmação para criar diretório de backup
2. Cria diretório com timestamp: `YYYYMMDDhhmmss_backup`
3. Para cada grupo de duplicatas:
   - Mantém o primeiro arquivo no local original
   - Pergunta sobre cada arquivo adicional
   - Move arquivos confirmados para backup, preservando estrutura

### 3. Estrutura de Backup
Se um arquivo duplicado está em `/home/user/docs/report.pdf` e o backup é `20250625124530_backup`, o novo local será:
```
./20250625124530_backup/home/user/docs/report.pdf
```

## Estrutura do Projeto

```
redup/
├── cmd/redup/          # Executável principal
│   └── main.go
├── scanner/            # Navegação de diretórios
│   ├── scanner.go
│   └── scanner_test.go
├── hasher/             # Cálculo de checksums
│   ├── hasher.go
│   └── hasher_test.go
├── deduplicator/       # Agrupamento e filtragem
│   ├── deduplicator.go
│   └── deduplicator_test.go
├── reporter/           # Relatórios e exportação
│   └── reporter.go
├── menu/               # Interface interativa
│   └── menu.go
├── backup/             # Gerenciamento de backup
│   └── backup.go
├── go.mod
├── go.sum
└── README.md
```

### Relação com o Projeto scopy

O `redup` reutiliza padrões de código do projeto [scopy](https://github.com/dakoctba/scopy), incluindo:
- **Arquitetura Modular**: Separação clara de responsabilidades em pacotes
- **Padrões de CLI**: Uso de flags e parsing de argumentos
- **Tratamento de Erros**: Estratégias consistentes de error handling
- **Testes Unitários**: Cobertura de testes para funcionalidades críticas

## Exemplos de Uso

### Exemplo 1: Limpeza de Downloads
```bash
# Encontrar duplicatas maiores que 1MB
redup --dir ~/Downloads --min-size 1048576

# Resultado:
# Found 3 duplicate groups.
# Total space that can be freed: 45.2 MB
```

### Exemplo 2: Backup de Fotos
```bash
# Usar MD5 para fotos (mais rápido)
redup --dir ~/Pictures --checksum md5 --backup-dir ~/photo_backups
```

### Exemplo 3: Análise sem Ação
```bash
# Apenas analisar sem mover arquivos
redup --dir ~/Documents --dry-run --json > duplicates.json
```

## Códigos de Saída

- `0` - Sucesso
- `1` - Cancelado pelo usuário
- `2` - Erro (diretório inexistente, permissões, etc.)

## Testes

Execute os testes unitários:
```bash
go test ./...
```

Execute testes com cobertura:
```bash
go test -cover ./...
```

## Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## Agradecimentos

- Inspirado no projeto [scopy](https://github.com/dakoctba/scopy)
- Comunidade Go por ferramentas e bibliotecas excelentes
