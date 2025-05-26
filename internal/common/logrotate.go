package common

import (
	"fmt"
	"os"
)

// RotateLog rotates log files up to maxBackups
func RotateLog(logFile string, maxBackups int) error {
	// Nothing to do if main log doesn't exist
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return nil
	}

	// Step 1: Shift old backups (.1 → .2 → ... → .5)
	for i := maxBackups - 1; i > 0; i-- {
		src := fmt.Sprintf("%s.%d", logFile, i)
		dst := fmt.Sprintf("%s.%d", logFile, i+1)

		// If target exists, remove it before renaming
		os.Remove(dst)

		// If source exists, move it
		if _, err := os.Stat(src); !os.IsNotExist(err) {
			if err := os.Rename(src, dst); err != nil {
				return fmt.Errorf("failed to rotate %s to %s: %w", src, dst, err)
			}
		}
	}

	// Step 2: Move current log to .1
	if err := os.Rename(logFile, fmt.Sprintf("%s.1", logFile)); err != nil {
		return fmt.Errorf("failed to rotate current log to .1: %w", err)
	}

	return nil
}
