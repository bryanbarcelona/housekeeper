package purge

import (
	"os"

	"github.com/spf13/afero"
)

var (
	// Make filesystem operations mockable
	osStat  = os.Stat
	osChmod = os.Chmod
	AppFs   = afero.NewOsFs()

	// Keep UnlockPath as a function but make it mockable
	UnlockPath = func(path string) error {
		info, err := osStat(path)
		if err != nil {
			return err
		}

		targetMode := os.FileMode(0666)
		if info.IsDir() {
			targetMode = 0777
		}

		if info.Mode().Perm() == targetMode {
			return nil
		}

		return osChmod(path, targetMode)
	}
)

/*
// UnlockPath makes a filesystem path (file or directory) writable for user operations.
// For files: Sets 0666 permissions (rw-rw-rw-)
// For directories: Sets 0777 permissions (rwxrwxrwx) to maintain traversability
func UnlockPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	targetMode := os.FileMode(0666) // Default for files
	if info.IsDir() {
		targetMode = 0777 // Directory needs execute bits
	}

	// Skip if already has correct permissions
	if info.Mode().Perm() == targetMode {
		return nil
	}

	return os.Chmod(path, targetMode)
}
*/
