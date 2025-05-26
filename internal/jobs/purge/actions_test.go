package purge

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestApply(t *testing.T) {
	// Setup common test variables
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "testfile.txt")
	renamedFile := filepath.Join(testDir, "renamed.txt")
	testDirToRemove := filepath.Join(testDir, "emptydir")

	// Create test file and directory
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := os.Mkdir(testDirToRemove, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name    string
		change  Change
		setup   func() error
		wantErr bool
		verify  func() error
	}{
		{
			name: "successful file deletion",
			change: Change{
				Type:   DeleteFile,
				Target: testFile,
			},
			verify: func() error {
				if _, err := os.Stat(testFile); !os.IsNotExist(err) {
					return errors.New("file still exists after deletion")
				}
				return nil
			},
		},
		{
			name: "file deletion fails when file doesn't exist",
			change: Change{
				Type:   DeleteFile,
				Target: "nonexistent.txt",
			},
			wantErr: true,
		},
		{
			name: "successful file rename",
			change: Change{
				Type:    RenameFile,
				Target:  testFile,
				NewName: renamedFile,
			},
			setup: func() error {
				// Recreate the test file since previous test might have deleted it
				return os.WriteFile(testFile, []byte("test content"), 0644)
			},
			verify: func() error {
				if _, err := os.Stat(testFile); !os.IsNotExist(err) {
					return errors.New("original file still exists after rename")
				}
				if _, err := os.Stat(renamedFile); err != nil {
					return errors.New("renamed file does not exist")
				}
				return nil
			},
		},
		{
			name: "file rename fails when source doesn't exist",
			change: Change{
				Type:    RenameFile,
				Target:  "nonexistent.txt",
				NewName: renamedFile,
			},
			wantErr: true,
		},
		{
			name: "successful directory removal",
			change: Change{
				Type:   RemoveDir,
				Target: testDirToRemove,
			},
			verify: func() error {
				if _, err := os.Stat(testDirToRemove); !os.IsNotExist(err) {
					return errors.New("directory still exists after removal")
				}
				return nil
			},
		},
		{
			name: "directory removal fails when not empty",
			change: Change{
				Type:   RemoveDir,
				Target: testDir,
			},
			wantErr: true,
		},
		{
			name: "unknown change type",
			change: Change{
				Type:   "invalid_type",
				Target: testFile,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup if needed
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Run the function
			err := Apply(tt.change)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify results if no error expected
			if !tt.wantErr && tt.verify != nil {
				if err := tt.verify(); err != nil {
					t.Errorf("Verification failed: %v", err)
				}
			}
		})
	}
}
