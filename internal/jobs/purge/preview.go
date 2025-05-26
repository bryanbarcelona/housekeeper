package purge

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

// RunDry performs housekeeping checks but does not modify anything.
func PreviewChanges(directory string, cfg *Config) ([]Change, error) {
	var changes []Change

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("[ERROR] Accessing %s: %v\n", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}

		// 1. Check if file should be deleted
		if c := checkDelete(path, cfg.ExtensionsToDelete); c != nil {
			changes = append(changes, *c)
		} else {
			// 2. If not deleting, try renaming (replacement > lowercase)
			if c := computeRename(path, cfg.ExtensionReplacements); c != nil {
				changes = append(changes, *c)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	emptyDirs, err := findEmptyDirs(directory, changes)
	if err != nil {
		return nil, err
	}
	changes = append(changes, emptyDirs...)

	return changes, nil
}

