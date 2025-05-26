package purge

import (
	"strings"
	"testing"
)

func TestCheckDelete(t *testing.T) {
	// Load production config using the same method as production code
	cfg, err := LoadConfigWithOptions(LoadConfigOptions{})
	if err != nil {
		t.Fatalf("Failed to load production config: %v", err)
	}

	// Core test cases for special files (not extension-dependent)
	coreTests := []struct {
		name string
		path string
		want *Change
	}{
		{
			name: "should delete dot underscore prefix files",
			path: "/path/to/._file",
			want: &Change{Type: DeleteFile, Target: "/path/to/._file"},
		},
		{
			name: "should delete .DS_Store files",
			path: "/path/to/.DS_Store",
			want: &Change{Type: DeleteFile, Target: "/path/to/.DS_Store"},
		},
		{
			name: "should delete hidden files",
			path: "/path/to/.hidden",
			want: &Change{Type: DeleteFile, Target: "/path/to/.hidden"},
		},
		{
			name: "should not delete normal files without matching extensions",
			path: "/path/to/normal_file",
			want: nil,
		},
	}

	// Run core tests
	for _, tt := range coreTests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkDelete(tt.path, cfg.ExtensionsToDelete)
			assertChangeEquals(t, got, tt.want)
		})
	}

	// Test every extension in the production config
	for _, ext := range cfg.ExtensionsToDelete {
		ext := ext // capture range variable

		t.Run("should delete files with "+ext+" extension", func(t *testing.T) {
			path := "/path/to/file" + ext
			got := checkDelete(path, cfg.ExtensionsToDelete)
			want := &Change{Type: DeleteFile, Target: path}
			assertChangeEquals(t, got, want)
		})

		t.Run("should delete files with uppercase "+ext+" extension", func(t *testing.T) {
			path := "/path/to/file" + strings.ToUpper(ext)
			got := checkDelete(path, cfg.ExtensionsToDelete)
			want := &Change{Type: DeleteFile, Target: path}
			assertChangeEquals(t, got, want)
		})
	}
}

func assertChangeEquals(t *testing.T, got, want *Change) {
	t.Helper()

	switch {
	case want == nil && got != nil:
		t.Fatalf("Expected nil, got %+v", got)
	case want != nil && got == nil:
		t.Fatalf("Expected %+v, got nil", want)
	case got != nil:
		if got.Type != want.Type {
			t.Errorf("Type mismatch: got %q, want %q", got.Type, want.Type)
		}
		if got.Target != want.Target {
			t.Errorf("Target mismatch: got %q, want %q", got.Target, want.Target)
		}
		if got.NewName != "" {
			t.Errorf("NewName should be empty, got %q", got.NewName)
		}
	}
}
