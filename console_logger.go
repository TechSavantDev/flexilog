package flexilog

import (
	"io"
	"os"
)

type consoleLogger struct {
	*logger
}

func ConsoleLogger(level LogLevel) Logger {
	return &consoleLogger{
		logger: &logger{
			level: level,
			out:   []io.WriteCloser{os.Stdout},
		},
	}
}
