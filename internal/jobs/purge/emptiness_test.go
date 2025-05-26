package purge

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

// setupEmptyDir creates a single empty directory for testing.
func setupEmptyDir(t *testing.T) (string, error) {
	t.Helper()
	testRoot := t.TempDir()
	dir := filepath.Join(testRoot, "empty_dir")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return testRoot, nil
}

// setupAllDeleteDir creates a directory with files all marked for deletion.
func setupAllDeleteDir(t *testing.T) (string, error) {
	t.Helper()
	testRoot := t.TempDir()
	dir := filepath.Join(testRoot, "all_delete_dir")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	files := []string{
		filepath.Join(dir, "delete1.tmp"),
		filepath.Join(dir, "delete2.txt"),
	}
	for _, file := range files {
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			return "", err
		}
	}
	return testRoot, nil
}

// setupMixedDir creates a directory with mixed content (empty and non-empty subdirs).
func setupMixedDir(t *testing.T) (string, error) {
	t.Helper()
	testRoot := t.TempDir()
	dirs := []string{
		filepath.Join(testRoot, "mixed_dir", "empty_subdir"),
		filepath.Join(testRoot, "mixed_dir", "non_empty_subdir"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}
	files := []string{
		filepath.Join(testRoot, "mixed_dir", "non_empty_subdir", "keep.txt"),
		filepath.Join(testRoot, "mixed_dir", "non_empty_subdir", "delete.tmp"),
	}
	for _, file := range files {
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			return "", err
		}
	}
	return testRoot, nil
}

// setupNestedEmptyDir creates a nested empty directory structure.
func setupNestedEmptyDir(t *testing.T) (string, error) {
	t.Helper()
	testRoot := t.TempDir()
	dirs := []string{
		filepath.Join(testRoot, "nested_empty_dir"),
		filepath.Join(testRoot, "nested_empty_dir", "subdir_a"),
		filepath.Join(testRoot, "nested_empty_dir", "subdir_b"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}
	return testRoot, nil
}

// TestFindEmptyDirs tests the findEmptyDirs function.
func TestFindEmptyDirs(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) (string, error)
		changes  []Change
		expected []Change
		wantErr  bool
	}{
		{
			name:    "EmptyDirectory",
			setup:   setupEmptyDir,
			changes: []Change{},
			expected: []Change{
				{Type: RemoveDir, Target: "."},
				{Type: RemoveDir, Target: "empty_dir"},
			},
			wantErr: false,
		},
		{
			name:  "AllFilesMarkedForDeletion",
			setup: setupAllDeleteDir,
			changes: []Change{
				{Type: DeleteFile, Target: filepath.Join("all_delete_dir", "delete1.tmp")},
				{Type: DeleteFile, Target: filepath.Join("all_delete_dir", "delete2.txt")},
			},
			expected: []Change{
				{Type: RemoveDir, Target: "."},
				{Type: RemoveDir, Target: "all_delete_dir"},
			},
			wantErr: false,
		},
		{
			name:  "NonEmptyDirectory",
			setup: setupMixedDir,
			changes: []Change{
				{Type: DeleteFile, Target: filepath.Join("mixed_dir", "non_empty_subdir", "delete.tmp")},
			},
			expected: []Change{
				{Type: RemoveDir, Target: filepath.Join("mixed_dir", "empty_subdir")},
			},
			wantErr: false,
		},
		{
			name:    "NestedEmptyDirectories",
			setup:   setupNestedEmptyDir,
			changes: []Change{},
			expected: []Change{
				{Type: RemoveDir, Target: "."},
				{Type: RemoveDir, Target: "nested_empty_dir"},
				{Type: RemoveDir, Target: filepath.Join("nested_empty_dir", "subdir_a")},
				{Type: RemoveDir, Target: filepath.Join("nested_empty_dir", "subdir_b")},
			},
			wantErr: false,
		},
		{
			name:  "MixedContentDirectory",
			setup: setupMixedDir,
			changes: []Change{
				{Type: DeleteFile, Target: filepath.Join("mixed_dir", "non_empty_subdir", "delete.tmp")},
			},
			expected: []Change{
				{Type: RemoveDir, Target: filepath.Join("mixed_dir", "empty_subdir")},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir, err := tt.setup(t)
			if err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			// Adjust changes to use absolute paths
			absChanges := make([]Change, len(tt.changes))
			for i, c := range tt.changes {
				absChanges[i] = Change{
					Type:    c.Type,
					Target:  filepath.Join(testDir, c.Target),
					NewName: c.NewName,
				}
			}

			// Run findEmptyDirs
			changes, err := findEmptyDirs(testDir, absChanges)
			if (err != nil) != tt.wantErr {
				t.Errorf("findEmptyDirs() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Convert paths to relative for comparison
			var relChanges []Change
			for _, c := range changes {
				relTarget, _ := filepath.Rel(testDir, c.Target)
				relChanges = append(relChanges, Change{Type: c.Type, Target: relTarget})
			}

			// Sort changes for consistent comparison
			sort.Slice(relChanges, func(i, j int) bool {
				if relChanges[i].Type != relChanges[j].Type {
					return relChanges[i].Type < relChanges[j].Type
				}
				return relChanges[i].Target < relChanges[j].Target
			})
			sort.Slice(tt.expected, func(i, j int) bool {
				if tt.expected[i].Type != tt.expected[j].Type {
					return tt.expected[i].Type < tt.expected[j].Type
				}
				return tt.expected[i].Target < tt.expected[j].Target
			})

			if !reflect.DeepEqual(relChanges, tt.expected) {
				t.Errorf("findEmptyDirs() = %v, want %v", relChanges, tt.expected)
			}
		})
	}
}

// TestFindEmptyDirsWithInaccessible tests findEmptyDirs with an inaccessible file.
func TestFindEmptyDirsWithInaccessible(t *testing.T) {
	testDir := t.TempDir()
	file := filepath.Join(testDir, "inaccessible.txt")
	if err := os.WriteFile(file, []byte("content"), 000); err != nil {
		t.Fatal(err)
	}

	changes := []Change{
		{Type: DeleteFile, Target: file},
	}

	// Run findEmptyDirs
	result, err := findEmptyDirs(testDir, changes)
	if err != nil {
		t.Errorf("findEmptyDirs() error = %v, want nil", err)
	}

	// Expect the directory to be marked for removal since the only file is marked for deletion
	expected := []Change{
		{Type: RemoveDir, Target: "."},
	}

	// Convert paths to relative
	var relChanges []Change
	for _, c := range result {
		relTarget, _ := filepath.Rel(testDir, c.Target)
		relChanges = append(relChanges, Change{Type: c.Type, Target: relTarget})
	}

	// Sort changes
	sort.Slice(relChanges, func(i, j int) bool {
		return relChanges[i].Target < relChanges[j].Target
	})

	if !reflect.DeepEqual(relChanges, expected) {
		t.Errorf("findEmptyDirs() = %v, want %v", relChanges, expected)
	}
}
