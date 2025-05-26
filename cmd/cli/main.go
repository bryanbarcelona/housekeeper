package main

import (
	"flag"
	"fmt"
	"log"

	"housekeeper/internal/common"
	"housekeeper/internal/jobs/purge"
)

func main() {
	dir := flag.String("dir", ".", "Directory to scan")
	exts := flag.String("exts", "userconfigs/extensions_to_delete.json", "Delete config")
	repls := flag.String("repls", "userconfigs/extension_replacements.json", "Replace config")
	apply := flag.Bool("apply", false, "Apply changes (default is dry run)")
	logToFile := flag.Bool("log-to-file", true, "Enable file-based logging")
	logPath := flag.String("log-path", "logs/toolkit.log", "Path to log file")
	debugLogs := flag.Bool("debug", false, "Enable debug-level logs")
	alsoPrint := flag.Bool("also-print-to-console", true, "Also print logs to console when logging to file")
	flag.Parse()

	common.SetupLogging(common.LoggingConfig{
		LogToFile:          *logToFile,
		LogFilePath:        *logPath,
		Debug:              *debugLogs,
		AlsoPrintToConsole: *alsoPrint, // if logging to file, default no; or set manually
	})

	cfg, err := purge.LoadConfigWithOptions(purge.LoadConfigOptions{
		DeleteConfigPath:  *exts,
		ReplaceConfigPath: *repls,
	})
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	job := purge.NewJob(*dir, cfg)
	changes, err := job.Plan()
	if err != nil {
		log.Fatalf("Error during plan: %v", err)
	}

	fmt.Printf("Found %d changes:\n", len(changes))
	for i, change := range changes {
		printChange(i+1, change)
	}

	if !*apply {
		fmt.Println("\nNo changes applied (use -apply to execute)")
		return
	}

	applied, err := purge.ApplyAll(changes)
	if err != nil {
		log.Fatalf("Error during apply: %v", err)
	}

	fmt.Printf("\nApplied %d changes:\n", len(applied))
	for i, change := range applied {
		fmt.Printf("%2d. [APPLIED] %s\n", i+1, describeChange(change))
	}
}

func printChange(i int, change purge.Change) {
	switch change.Type {
	case purge.DeleteFile:
		fmt.Printf("%2d. [DELETE]     %s\n", i, change.Target)
	case purge.RenameFile:
		fmt.Printf("%2d. [RENAME]     %s → %s\n", i, change.Target, change.NewName)
	case purge.RemoveDir:
		fmt.Printf("%2d. [REMOVE DIR] %s\n", i, change.Target)
	default:
		fmt.Printf("%2d. [UNKNOWN]    %s (%s)\n", i, change.Target, change.Type)
	}
}

func describeChange(change purge.Change) string {
	switch change.Type {
	case purge.DeleteFile:
		return fmt.Sprintf("Deleted %s", change.Target)
	case purge.RenameFile:
		return fmt.Sprintf("Renamed %s → %s", change.Target, change.NewName)
	case purge.RemoveDir:
		return fmt.Sprintf("Removed empty dir %s", change.Target)
	default:
		return fmt.Sprintf("Unknown action on %s", change.Target)
	}
}
