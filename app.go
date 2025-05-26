package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
    "bytes"
)

// App struct
type App struct {
    ctx context.Context
}

// Config holds the extensions and prefixes to delete
type Config struct {
    ExtensionsToDelete []string `json:"extensions_to_delete"`
    PrefixesToDelete   []string `json:"prefixes_to_delete"`
}

// NewApp creates a new App application struct
func NewApp() *App {
    return &App{}
}

// startup is called when the app starts. The context is saved
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

// LoadConfig loads JSON config from a file path
func LoadConfig(configPath string) (*Config, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }
    var cfg Config
    err = json.Unmarshal(data, &cfg)
    return &cfg, err
}

// ShouldDelete decides if a file should be deleted based on config
func ShouldDelete(fileName string, cfg *Config) bool {
    nameLower := strings.ToLower(fileName)
    for _, prefix := range cfg.PrefixesToDelete {
        if strings.HasPrefix(nameLower, strings.ToLower(prefix)) {
            return true
        }
    }
    for _, ext := range cfg.ExtensionsToDelete {
        if strings.HasSuffix(nameLower, strings.ToLower(ext)) {
            return true
        }
    }
    return false
}

// SimulateScan scans the directory and returns the simulation results as a string
func (a *App) SimulateScan(root string) string {
    configPath := "config.json" // or pass as parameter if you want

    cfg, err := LoadConfig(configPath)
    if err != nil {
        return fmt.Sprintf("Failed to load config: %v", err)
    }

    var buffer bytes.Buffer

    buffer.WriteString(fmt.Sprintf("Scanning directory: %s\n", root))

    err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            buffer.WriteString(fmt.Sprintf("Error accessing path %s: %v\n", path, err))
            return nil
        }

        if !d.IsDir() {
            buffer.WriteString(fmt.Sprintf("Found: %s\n", path))
            if ShouldDelete(d.Name(), cfg) {
                buffer.WriteString(fmt.Sprintf("  --> Would delete: %s\n", path))
            }
        }
        return nil
    })

    if err != nil {
        buffer.WriteString(fmt.Sprintf("Error walking directory: %v\n", err))
    }

    return buffer.String()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
    return fmt.Sprintf("Hello %s, It's show time!", name)
}
