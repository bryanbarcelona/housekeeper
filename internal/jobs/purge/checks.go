package purge

import (
    "path/filepath"
    "strings"
)

func checkDelete(path string, extensions []string) *Change {
    name := filepath.Base(path)
    lower := strings.ToLower(name)

    for _, ext := range extensions {
        if strings.HasSuffix(lower, ext) {
            return &Change{Type: DeleteFile, Target: path}
        }
    }

    if strings.HasPrefix(name, "._") || strings.HasPrefix(name, ".DS_Store") || strings.HasPrefix(name, ".") {
        return &Change{Type: DeleteFile, Target: path}
    }

    return nil
}
