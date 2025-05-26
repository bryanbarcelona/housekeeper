package purge

import (
	"fmt"
)

// ApplyAll applies all given changes
func ApplyAll(changes []Change) ([]Change, error) {
	var applied []Change
	var applyErr error

	for _, change := range changes {
		if err := Apply(change); err != nil {
			fmt.Printf("[ERROR] Failed to apply change: %v\n", err)
			if applyErr == nil {
				applyErr = err
			}
			continue
		}
		applied = append(applied, change)
	}

	return applied, applyErr
}
