package purge

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestApplyAll tests the ApplyAll function with various scenarios
func TestApplyAll(t *testing.T) {
	// Setup mock filesystem
	fs := afero.NewMemMapFs()
	AppFs = fs // Assuming appFs is the package-level filesystem variable

	// Create some test files and directories
	testDir := filepath.Join("test", "dir")
	testFile1 := filepath.Join("test", "file1.txt")
	testFile2 := filepath.Join("test", "file2.txt")
	require.NoError(t, fs.MkdirAll(testDir, 0755))
	require.NoError(t, afero.WriteFile(fs, testFile1, []byte("content"), 0644))
	require.NoError(t, afero.WriteFile(fs, testFile2, []byte("content"), 0644))

	// Mock UnlockPath function for testing
	originalUnlockPath := UnlockPath
	defer func() { UnlockPath = originalUnlockPath }()
	UnlockPath = func(path string) error { return nil }

	tests := []struct {
		name        string
		changes     []Change
		wantApplied int
		wantErr     bool
		setup       func()
		verify      func(t *testing.T)
	}{
		{
			name: "successful file deletion",
			changes: []Change{
				{Type: DeleteFile, Target: testFile1},
			},
			wantApplied: 1,
			verify: func(t *testing.T) {
				exists, err := afero.Exists(fs, testFile1)
				assert.NoError(t, err)
				assert.False(t, exists, "file should be deleted")
			},
		},
		{
			name: "successful file rename",
			changes: []Change{
				{
					Type:    RenameFile,
					Target:  testFile1,
					NewName: filepath.Join("test", "renamed.txt"),
				},
			},
			setup: func() {
				// Ensure the file exists before rename
				afero.WriteFile(fs, testFile1, []byte("content"), 0644)
			},
			wantApplied: 1,
			verify: func(t *testing.T) {
				// Check original doesn't exist
				exists, err := afero.Exists(fs, testFile1)
				assert.NoError(t, err)
				assert.False(t, exists, "original file should not exist")

				// Check new file exists
				renamedPath := filepath.Join("test", "renamed.txt")
				exists, err = afero.Exists(fs, renamedPath)
				assert.NoError(t, err)
				assert.True(t, exists, "renamed file should exist")
			},
		},
		{
			name: "successful directory removal",
			setup: func() {
				emptyDir := filepath.Join("empty", "dir")
				require.NoError(t, fs.MkdirAll(emptyDir, 0755))
			},
			changes: []Change{
				{Type: RemoveDir, Target: filepath.Join("empty", "dir")},
			},
			wantApplied: 1,
			verify: func(t *testing.T) {
				exists, err := afero.Exists(fs, filepath.Join("empty", "dir"))
				assert.NoError(t, err)
				assert.False(t, exists, "directory should be removed")
			},
		},
		{
			name: "mixed successful and failed changes",
			setup: func() {
				// Recreate files that might have been deleted by previous tests
				afero.WriteFile(fs, testFile1, []byte("content"), 0644)
				afero.WriteFile(fs, testFile2, []byte("content"), 0644)
			},
			changes: []Change{
				{Type: DeleteFile, Target: testFile1},
				{Type: DeleteFile, Target: "nonexistent.txt"},
				{Type: DeleteFile, Target: testFile2},
			},
			wantApplied: 2,
			wantErr:     true,
			verify: func(t *testing.T) {
				for _, file := range []string{testFile1, testFile2} {
					exists, err := afero.Exists(fs, file)
					assert.NoError(t, err)
					assert.False(t, exists, "file should be deleted")
				}
			},
		},
		{
			name: "unknown change type",
			changes: []Change{
				{Type: "invalid_type", Target: testFile1},
			},
			wantApplied: 0,
			wantErr:     true,
		},
		{
			name:        "empty changes",
			changes:     []Change{},
			wantApplied: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run setup if provided
			if tt.setup != nil {
				tt.setup()
			}

			// Execute ApplyAll
			applied, err := ApplyAll(tt.changes)

			// Check error expectation
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Check number of applied changes
			assert.Equal(t, tt.wantApplied, len(applied), "unexpected number of applied changes")

			// Run verification if provided
			if tt.verify != nil {
				tt.verify(t)
			}
		})
	}
}
