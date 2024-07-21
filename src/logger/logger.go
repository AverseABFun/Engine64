package logger

import (
	"log"
	"os"

	formatColor "github.com/fatih/color"
)

const (
	// LogInfo is used to log information messages
	LogInfo = iota
	// LogError is used to log error messages
	LogError
	// LogWarning is used to log warning messages
	LogWarning
	// LogDebug is used to log debug messages
	LogDebug
)

var logger = log.New(os.Stderr, "", log.Lmsgprefix|log.Ltime)

func Log(msg string, logType int) {
	msg = formatColor.MagentaString(msg)
	switch logType {
	case LogInfo:
		logger.SetPrefix(formatColor.GreenString("[INFO] "))
	case LogError:
		logger.SetPrefix(formatColor.RedString("[ERROR] "))
	case LogWarning:
		logger.SetPrefix(formatColor.YellowString("[WARNING] "))
	case LogDebug:
		logger.SetPrefix(formatColor.BlueString("[DEBUG] "))
	default:
		logger.SetPrefix(formatColor.WhiteString("[INFO] "))
	}
	logger.Println(msg)
}

func LogEmptyNewline() {
	logger.SetPrefix("")
	logger.SetFlags(0)
	logger.Print("\n")
	logger.SetFlags(log.Lmsgprefix | log.Ltime)
}
