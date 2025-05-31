package main

import (
    "context"
    "housekeeper/internal/common"
    "housekeeper/internal/jobs/purge"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
    ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
    return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Set up logging (same as CLI)
    common.SetupLogging(common.LoggingConfig{
        LogToFile:          true,
        LogFilePath:        "logs/toolkit.log",
        Debug:              false,
        AlsoPrintToConsole: true,
    })
}

// GetChanges runs the purge job and returns the list of changes
func (a *App) GetChanges(dir string, deleteConfigPath string, replaceConfigPath string) ([]Change, error) {
    // Load configuration
    cfg, err := purge.LoadConfigWithOptions(purge.LoadConfigOptions{
        DeleteConfigPath:  deleteConfigPath,
        ReplaceConfigPath: replaceConfigPath,
    })
    if err != nil {
        return nil, err
    }

    // Create and run the purge job
    job := purge.NewJob(dir, cfg)
    changes, err := job.Plan()
    if err != nil {
        return nil, err
    }

    // Convert purge.Change to a JSON-serializable struct
    result := make([]Change, len(changes))
    for i, change := range changes {
        result[i] = Change{
            Type:     string(change.Type),
            Target:   change.Target,
            NewName:  change.NewName,
            Selected: true, // Default: checked
        }
    }

    return result, nil
}

// OpenDirectoryDialog opens a directory selection dialog
func (a *App) OpenDirectoryDialog(title string, defaultDirectory string) (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            title,
		DefaultDirectory: defaultDirectory,
	})
}

// Change struct for JSON serialization
type Change struct {
    Type     string `json:"type"`
    Target   string `json:"target"`
    NewName  string `json:"newName"`
    Selected bool   `json:"selected"`
}