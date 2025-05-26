package purge

import (
	"fmt"

	"housekeeper/internal/common"
)

func Apply(change Change) error {
	switch change.Type {
	case DeleteFile:
		common.Info.Printf("Deleting %s\n", change.Target)
		if err := UnlockPath(change.Target); err != nil {
			common.Warn.Printf("Unlock failed: %v", err)
			return fmt.Errorf("unlocking %s: %w", change.Target, err)
		}
		if err := AppFs.Remove(change.Target); err != nil {
			common.Error.Printf("Failed to delete %s: %v", change.Target, err)
			return fmt.Errorf("deleting %s: %w", change.Target, err)
		}
	case RenameFile:
		common.Info.Printf("Renaming %s → %s\n", change.Target, change.NewName)
		if err := UnlockPath(change.Target); err != nil {
			common.Warn.Printf("Unlock failed: %v", err)
			return fmt.Errorf("unlocking %s: %w", change.Target, err)
		}
		if err := AppFs.Rename(change.Target, change.NewName); err != nil {
			common.Error.Printf("Failed to rename %s: %v", change.Target, err)
			return fmt.Errorf("renaming %s → %s: %w", change.Target, change.NewName, err)
		}
	case RemoveDir:
		common.Info.Printf("Removing empty directory %s\n", change.Target)
		if err := AppFs.Remove(change.Target); err != nil { // Use RemoveAll instead of Remove
			common.Error.Printf("Failed to remove dir %s: %v", change.Target, err)
			return fmt.Errorf("removing dir %s: %w", change.Target, err)
		}
	default:
		common.Warn.Printf("Unknown change type: %v", change.Type)
		return fmt.Errorf("unknown change type: %v", change.Type)
	}
	return nil
}
