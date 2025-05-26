package purge

import (
    "path/filepath"
    "strings"
)

func computeRename(path string, replacements map[string]string) *Change {
    ext := strings.ToLower(filepath.Ext(path))
    base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

    // Try replacement first
    if newExt, ok := replacements[ext]; ok {
        newPath := filepath.Join(filepath.Dir(path), base+newExt)
        return &Change{Type: RenameFile, Target: path, NewName: newPath}
    }

    // Fallback to lowercase if no replacement
    originalExt := filepath.Ext(path)
    lowerExt := strings.ToLower(originalExt)
    if originalExt != lowerExt {
        newPath := filepath.Join(filepath.Dir(path), base+lowerExt)
        return &Change{Type: RenameFile, Target: path, NewName: newPath}
    }

    return nil
}