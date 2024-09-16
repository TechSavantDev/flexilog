package main

import "github.com/TechSavantDev/flexilog"

func main() {
	fileLogger, _ := flexilog.FileLogger(flexilog.INFO, "./info.log")
	fileLogger.Info("testing info logger")
	fileLogger.Debug("testing debug logger")
	fileLogger.Warn("testing warn logger")
	fileLogger.Error("testing error logger")

	consoleLogger := flexilog.ConsoleLogger(flexilog.DEBUG)
	consoleLogger.Info("testing console logger")
	consoleLogger.Warn("testing warn logger")
	consoleLogger.Error("testing error logger")
	consoleLogger.Debug("testing console logger")
}
