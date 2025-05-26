// purge/config.go
package purge

import (
    "encoding/json"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "runtime"
)

// LoadConfigWithOptions loads config using optional paths + fallback by default
func LoadConfigWithOptions(opts LoadConfigOptions) (*Config, error) {
    deletePath := opts.DeleteConfigPath
    replacePath := opts.ReplaceConfigPath

    // Always try fallback if paths missing
    if deletePath == "" || replacePath == "" {
        root, err := findProjectRoot()
        if err != nil {
            return nil, fmt.Errorf("unable to locate project root: %w", err)
        }

        if deletePath == "" {
            deletePath = filepath.Join(root, "userconfigs", "extensions_to_delete.json")
        }
        if replacePath == "" {
            replacePath = filepath.Join(root, "userconfigs", "extension_replacements.json")
        }
    }

    delData, err := readJSONFile(deletePath)
    if err != nil {
        return nil, fmt.Errorf("reading delete config: %w", err)
    }

    replData, err := readJSONFile(replacePath)
    if err != nil {
        return nil, fmt.Errorf("reading replace config: %w", err)
    }

    var exts []string
    if err := json.Unmarshal(delData, &exts); err != nil {
        return nil, fmt.Errorf("parsing delete config: %w", err)
    }

    var repls map[string]string
    if err := json.Unmarshal(replData, &repls); err != nil {
        return nil, fmt.Errorf("parsing replace config: %w", err)
    }

    return &Config{
        ExtensionsToDelete:    exts,
        ExtensionReplacements: repls,
    }, nil
}

func readJSONFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fs.ErrNotExist
        }
        return nil, err
    }
    return data, nil
}

// findProjectRoot finds the root of the project by looking upward until it finds userconfigs/
func findProjectRoot() (string, error) {
    _, filename, _, _ := runtime.Caller(0)
    dir := filepath.Dir(filename)

    for {
        if _, err := os.Stat(filepath.Join(dir, "userconfigs")); err == nil {
            return dir, nil
        }
        parent := filepath.Dir(dir)
        if parent == dir {
            return "", fmt.Errorf("could not find project root")
        }
        dir = parent
    }
}