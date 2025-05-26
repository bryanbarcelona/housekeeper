package purge

// ChangeType represents the type of change to apply
type ChangeType string

const (
    DeleteFile ChangeType = "delete_file"
    RenameFile ChangeType = "rename_file"
    RemoveDir  ChangeType = "remove_dir"
)

// Change describes a single planned or applied change
type Change struct {
    Type    ChangeType `json:"type"`
    Target  string     `json:"target"`
    NewName string     `json:"new_name,omitempty"` // only used for rename
}

// Config holds settings loaded from JSON files
type Config struct {
    ExtensionsToDelete    []string          `json:"extensions_to_delete"`
    ExtensionReplacements map[string]string `json:"extension_replacements"`
}

// LoadConfigOptions holds optional overrides for config paths
type LoadConfigOptions struct {
    DeleteConfigPath  string // Optional override
    ReplaceConfigPath string // Optional override
}