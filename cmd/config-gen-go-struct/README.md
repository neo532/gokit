# config-gen-go-struct

Generate Go struct code from YAML/JSON/INI config files.

## Usage

```bash
# Process all yaml files in a directory (default)
config-gen-go-struct ./configs

# Specify format
config-gen-go-struct -format json ./configs
config-gen-go-struct -format json,yaml ./configs

# Single file
config-gen-go-struct config.yaml

# Custom package name and root struct name
config-gen-go-struct -pkg config -type Config ./configs
```

## Output

Input `data.yaml` produces `data.cfg.go`. When processing multiple files, an additional `config.cfg.go` is generated as the aggregation entry point.

## Clean up generated files

```bash
# zsh
rm -f *.cfg.go(N)

# bash
rm -f *.cfg.go 2>/dev/null
```
