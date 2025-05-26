package purge

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUnlockPath(t *testing.T) {
	// Create a temporary test directory
	tempDir := t.TempDir()

	// Test cases
	tests := []struct {
		name        string
		setup       func() string // returns file path
		wantErr     bool
		wantMode    os.FileMode // expected permission mode
		checkResult func(path string) error
	}{
		{
			name: "successfully unlock read-only file",
			setup: func() string {
				path := filepath.Join(tempDir, "readonly.txt")
				err := os.WriteFile(path, []byte("test"), 0444)
				if err != nil {
					t.Fatal(err)
				}
				return path
			},
			wantErr:  false,
			wantMode: 0666,
		},
		{
			name: "already writable file - no change needed",
			setup: func() string {
				path := filepath.Join(tempDir, "writable.txt")
				err := os.WriteFile(path, []byte("test"), 0666)
				if err != nil {
					t.Fatal(err)
				}
				return path
			},
			wantErr:  false,
			wantMode: 0666,
		},
		{
			name: "non-existent file - should error",
			setup: func() string {
				return filepath.Join(tempDir, "nonexistent.txt")
			},
			wantErr: true,
		},
		{
			name: "unlock directory successfully",
			setup: func() string {
				path := filepath.Join(tempDir, "subdir")
				err := os.Mkdir(path, 0555)
				if err != nil {
					t.Fatal(err)
				}
				return path
			},
			wantErr:  false,
			wantMode: 0777,
		},
		{
			name: "unlock read-only directory",
			setup: func() string {
				path := filepath.Join(tempDir, "readonly_dir")
				err := os.Mkdir(path, 0555) // r-xr-xr-x
				if err != nil {
					t.Fatal(err)
				}
				return path
			},
			wantErr:  false,
			wantMode: 0777,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			err := UnlockPath(path)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnlockPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.wantMode != 0 {
				info, err := os.Stat(path)
				if err != nil {
					t.Errorf("Failed to stat path after UnlockPath: %v", err)
					return
				}
				if info.Mode().Perm() != tt.wantMode {
					t.Errorf("Permissions = %04o, want %04o", info.Mode().Perm(), tt.wantMode)
				}
			}
		})
	}
}
