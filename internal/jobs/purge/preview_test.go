package purge

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// createTestDir creates a temporary directory structure for testing.
func createTestDir(t *testing.T) (string, error) {
	t.Helper()

	// Use t.TempDir for automatic cleanup
	testRoot := t.TempDir()

	// Create our test directory structure
	dirs := []string{
		filepath.Join(testRoot, "nested_empty_dir"),
		filepath.Join(testRoot, "nested_empty_dir", "empty_subdir_a"),
		filepath.Join(testRoot, "nested_empty_dir", "empty_subdir_b"),
		filepath.Join(testRoot, "nested_empty_dir", "empty_subdir_b", "empty_subdir_b1"),
		filepath.Join(testRoot, "nested_mixed_content_dir"),
		filepath.Join(testRoot, "nested_mixed_content_dir", "subdir_w_keep_files"),
		filepath.Join(testRoot, "nested_mixed_content_dir", "subdir_w_mixed_content"),
		filepath.Join(testRoot, "nested_mixed_content_dir", "subdir_w_temp_files"),
		filepath.Join(testRoot, "top_lvl_empty_dir"),
		filepath.Join(testRoot, "top_lvl_mixed_temp_files_dir"),
		filepath.Join(testRoot, "top_lvl_only_temp_files_dir"),
	}

	// Create all directories
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create test files
	files := map[string]string{
		// Nested mixed content dir files
		filepath.Join(testRoot, "nested_mixed_content_dir", "keep_file_4.txt"):                        "content",
		filepath.Join(testRoot, "nested_mixed_content_dir", "subdir_w_keep_files", "keep_file_2.txt"): "content",

		// Subdir with mixed content
		filepath.Join(testRoot, "nested_mixed_content_dir", "subdir_w_mixed_content", ".dot_temp_file_2.txt"):     "content",
		filepath.Join(testRoot, "nested_mixed_content_dir", "subdir_w_mixed_content", "keep_file_3.txt"):          "content",
		filepath.Join(testRoot, "nested_mixed_content_dir", "subdir_w_mixed_content", "temp_file_w_suffix_2.tmp"): "content",

		// Subdir with temp files
		filepath.Join(testRoot, "nested_mixed_content_dir", "subdir_w_temp_files", ".dot_temp_file_3.txt"):     "content",
		filepath.Join(testRoot, "nested_mixed_content_dir", "subdir_w_temp_files", "temp_file_w_suffix_3.tmp"): "content",

		// Top level mixed files
		filepath.Join(testRoot, "top_lvl_mixed_temp_files_dir", ".dot_temp_file.txt"):     "content",
		filepath.Join(testRoot, "top_lvl_mixed_temp_files_dir", "keep_file.txt"):          "content",
		filepath.Join(testRoot, "top_lvl_mixed_temp_files_dir", "temp_file_w_suffix.tmp"): "content",

		// Top level only temp files
		filepath.Join(testRoot, "top_lvl_only_temp_files_dir", ".dot_temp_file_4.txt"):     "content",
		filepath.Join(testRoot, "top_lvl_only_temp_files_dir", "temp_file_w_suffix_4.tmp"): "content",
	}

	for file, content := range files {
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("failed to create file %s: %w", file, err)
		}
	}

	return testRoot, nil
}

// TestPreviewChanges tests the PreviewChanges function with a real filesystem.
func TestPreviewChanges(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		expected []Change
	}{
		{
			name: "DeleteTempAndHiddenFiles",
			cfg: &Config{
				ExtensionsToDelete:    []string{".tmp"},
				ExtensionReplacements: map[string]string{".txt": ".bak"},
			},
			expected: []Change{
				{Type: DeleteFile, Target: filepath.Join("nested_mixed_content_dir", "subdir_w_mixed_content", ".dot_temp_file_2.txt")},
				{Type: DeleteFile, Target: filepath.Join("nested_mixed_content_dir", "subdir_w_mixed_content", "temp_file_w_suffix_2.tmp")},
				{Type: DeleteFile, Target: filepath.Join("nested_mixed_content_dir", "subdir_w_temp_files", ".dot_temp_file_3.txt")},
				{Type: DeleteFile, Target: filepath.Join("nested_mixed_content_dir", "subdir_w_temp_files", "temp_file_w_suffix_3.tmp")},
				{Type: DeleteFile, Target: filepath.Join("top_lvl_mixed_temp_files_dir", ".dot_temp_file.txt")},
				{Type: DeleteFile, Target: filepath.Join("top_lvl_mixed_temp_files_dir", "temp_file_w_suffix.tmp")},
				{Type: DeleteFile, Target: filepath.Join("top_lvl_only_temp_files_dir", ".dot_temp_file_4.txt")},
				{Type: DeleteFile, Target: filepath.Join("top_lvl_only_temp_files_dir", "temp_file_w_suffix_4.tmp")},
				{Type: RenameFile, Target: filepath.Join("nested_mixed_content_dir", "keep_file_4.txt"), NewName: filepath.Join("nested_mixed_content_dir", "keep_file_4.bak")},
				{Type: RenameFile, Target: filepath.Join("nested_mixed_content_dir", "subdir_w_keep_files", "keep_file_2.txt"), NewName: filepath.Join("nested_mixed_content_dir", "subdir_w_keep_files", "keep_file_2.bak")},
				{Type: RenameFile, Target: filepath.Join("nested_mixed_content_dir", "subdir_w_mixed_content", "keep_file_3.txt"), NewName: filepath.Join("nested_mixed_content_dir", "subdir_w_mixed_content", "keep_file_3.bak")},
				{Type: RenameFile, Target: filepath.Join("top_lvl_mixed_temp_files_dir", "keep_file.txt"), NewName: filepath.Join("top_lvl_mixed_temp_files_dir", "keep_file.bak")},
				{Type: RemoveDir, Target: filepath.Join("nested_empty_dir", "empty_subdir_a")},
				{Type: RemoveDir, Target: filepath.Join("nested_empty_dir", "empty_subdir_b", "empty_subdir_b1")},
				{Type: RemoveDir, Target: filepath.Join("nested_empty_dir", "empty_subdir_b")},
				{Type: RemoveDir, Target: filepath.Join("nested_empty_dir")},
				{Type: RemoveDir, Target: filepath.Join("nested_mixed_content_dir", "subdir_w_temp_files")},
				{Type: RemoveDir, Target: filepath.Join("top_lvl_empty_dir")},
				{Type: RemoveDir, Target: filepath.Join("top_lvl_only_temp_files_dir")},
			},
		},
		{
			name: "NoChanges",
			cfg: &Config{
				ExtensionsToDelete:    []string{".unknown"},
				ExtensionReplacements: map[string]string{},
			},
			expected: []Change{
				{Type: DeleteFile, Target: filepath.Join("nested_mixed_content_dir", "subdir_w_mixed_content", ".dot_temp_file_2.txt")},
				{Type: DeleteFile, Target: filepath.Join("nested_mixed_content_dir", "subdir_w_temp_files", ".dot_temp_file_3.txt")},
				{Type: DeleteFile, Target: filepath.Join("top_lvl_mixed_temp_files_dir", ".dot_temp_file.txt")},
				{Type: DeleteFile, Target: filepath.Join("top_lvl_only_temp_files_dir", ".dot_temp_file_4.txt")},
				{Type: RemoveDir, Target: filepath.Join("nested_empty_dir", "empty_subdir_a")},
				{Type: RemoveDir, Target: filepath.Join("nested_empty_dir", "empty_subdir_b", "empty_subdir_b1")},
				{Type: RemoveDir, Target: filepath.Join("nested_empty_dir", "empty_subdir_b")},
				{Type: RemoveDir, Target: filepath.Join("nested_empty_dir")},
				{Type: RemoveDir, Target: filepath.Join("top_lvl_empty_dir")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir, err := createTestDir(t)
			if err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			changes, err := PreviewChanges(testDir, tt.cfg)
			if err != nil {
				t.Fatalf("PreviewChanges failed: %v", err)
			}

			// Convert paths to relative for comparison
			var relChanges []Change
			for _, c := range changes {
				relTarget, _ := filepath.Rel(testDir, c.Target)
				relNewName := c.NewName
				if c.NewName != "" {
					relNewName, _ = filepath.Rel(testDir, c.NewName)
				}
				relChanges = append(relChanges, Change{Type: c.Type, Target: relTarget, NewName: relNewName})
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
				t.Errorf("PreviewChanges() = %v, want %v", relChanges, tt.expected)
			}
		})
	}
}

// previewChangesWithDeps is a test-specific version of PreviewChanges with injectable dependencies.
func previewChangesWithDeps(
	directory string,
	cfg *Config,
	checkDeleteFn func(path string, extensions []string) *Change,
	computeRenameFn func(path string, replacements map[string]string) *Change,
	findEmptyDirsFn func(root string, changes []Change) ([]Change, error),
	walkDirFn func(root string, fn fs.WalkDirFunc) error,
) ([]Change, error) {
	var changes []Change
	err := walkDirFn(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("[ERROR] Accessing %s: %v\n", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if c := checkDeleteFn(path, cfg.ExtensionsToDelete); c != nil {
			changes = append(changes, *c)
		} else if c := computeRenameFn(path, cfg.ExtensionReplacements); c != nil {
			changes = append(changes, *c)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	emptyDirs, err := findEmptyDirsFn(directory, changes)
	if err != nil {
		return nil, err
	}
	changes = append(changes, emptyDirs...)
	return changes, nil
}

// TestPreviewChangesUnit tests PreviewChanges with mocked dependencies.
func TestPreviewChangesUnit(t *testing.T) {
	tests := []struct {
		name            string
		dir             string
		cfg             *Config
		checkDeleteFn   func(path string, extensions []string) *Change
		computeRenameFn func(path string, replacements map[string]string) *Change
		findEmptyDirsFn func(root string, changes []Change) ([]Change, error)
		walkDirFn       func(root string, fn fs.WalkDirFunc) error
		expected        []Change
		wantErr         bool
	}{
		{
			name: "BasicFlow",
			dir:  "/test",
			cfg: &Config{
				ExtensionsToDelete:    []string{".tmp"},
				ExtensionReplacements: map[string]string{".txt": ".bak"},
			},
			checkDeleteFn: func(path string, extensions []string) *Change {
				if strings.HasSuffix(path, ".tmp") {
					return &Change{Type: DeleteFile, Target: path}
				}
				return nil
			},
			computeRenameFn: func(path string, replacements map[string]string) *Change {
				if strings.HasSuffix(path, ".txt") {
					return &Change{Type: RenameFile, Target: path, NewName: strings.TrimSuffix(path, ".txt") + ".bak"}
				}
				return nil
			},
			findEmptyDirsFn: func(root string, changes []Change) ([]Change, error) {
				return []Change{{Type: RemoveDir, Target: filepath.ToSlash(filepath.Join(root, "empty"))}}, nil
			},
			walkDirFn: func(root string, fn fs.WalkDirFunc) error {
				// Simulate walking two files
				fn("/test/file.tmp", &mockDirEntry{isDir: false}, nil)
				fn("/test/file.txt", &mockDirEntry{isDir: false}, nil)
				return nil
			},
			expected: []Change{
				{Type: DeleteFile, Target: "/test/file.tmp"},
				{Type: RenameFile, Target: "/test/file.txt", NewName: "/test/file.bak"},
				{Type: RemoveDir, Target: "/test/empty"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes, err := previewChangesWithDeps(
				tt.dir,
				tt.cfg,
				tt.checkDeleteFn,
				tt.computeRenameFn,
				tt.findEmptyDirsFn,
				tt.walkDirFn,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("previewChangesWithDeps() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Normalize paths for comparison
			for i, c := range changes {
				changes[i].Target = filepath.ToSlash(c.Target)
				if c.NewName != "" {
					changes[i].NewName = filepath.ToSlash(c.NewName)
				}
			}

			if !reflect.DeepEqual(changes, tt.expected) {
				t.Errorf("previewChangesWithDeps() = %v, want %v", changes, tt.expected)
			}
		})
	}
}

// mockDirEntry for filepath.WalkDir
type mockDirEntry struct {
	isDir bool
}

func (m *mockDirEntry) IsDir() bool                { return m.isDir }
func (m *mockDirEntry) Type() fs.FileMode          { return 0 }
func (m *mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }
func (m *mockDirEntry) Name() string               { return "" }

// TestPreviewChangesEdgeCases tests edge cases for PreviewChanges.
func TestPreviewChangesEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) (string, *Config)
		wantErr  bool
		expected []Change
	}{
		{
			name: "EmptyDirectory",
			setup: func(t *testing.T) (string, *Config) {
				dir := t.TempDir()
				return dir, &Config{ExtensionsToDelete: []string{".tmp"}}
			},
			wantErr: false,
			expected: []Change{
				{Type: RemoveDir, Target: "."},
			},
		},
		{
			name: "InaccessibleFile",
			setup: func(t *testing.T) (string, *Config) {
				dir := t.TempDir()
				file := filepath.Join(dir, "test.txt")
				if err := os.WriteFile(file, []byte("content"), 000); err != nil {
					t.Fatal(err)
				}
				return dir, &Config{ExtensionsToDelete: []string{".txt"}}
			},
			wantErr: false, // Errors are logged, not returned
			expected: []Change{
				{Type: DeleteFile, Target: "test.txt"},
				{Type: RemoveDir, Target: "."},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cfg := tt.setup(t)
			changes, err := PreviewChanges(dir, cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("PreviewChanges() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Convert paths to relative for comparison
			var relChanges []Change
			for _, c := range changes {
				relTarget, _ := filepath.Rel(dir, c.Target)
				relNewName := c.NewName
				if c.NewName != "" {
					relNewName, _ = filepath.Rel(dir, c.NewName)
				}
				relChanges = append(relChanges, Change{Type: c.Type, Target: relTarget, NewName: relNewName})
			}

			// Sort changes for consistent comparison
			sort.Slice(relChanges, func(i, j int) bool {
				if relChanges[i].Type != relChanges[j].Type {
					return relChanges[i].Type < relChanges[j].Type
				}
				return relChanges[i].Target < relChanges[j].Target
			})

			if !reflect.DeepEqual(relChanges, tt.expected) {
				t.Errorf("PreviewChanges() = %v, want %v", relChanges, tt.expected)
			}
		})
	}
}
