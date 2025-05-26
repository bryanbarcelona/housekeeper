package common

type LoggingConfig struct {
	LogToFile          bool
	LogFilePath        string
	Debug              bool
	AlsoPrintToConsole bool // ← new field
}
