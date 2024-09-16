package flexilog

import (
	"fmt"
	"io"
	"os"
)

type fileLogger struct {
	*logger
	file *os.File
}

func FileLogger(level LogLevel, path string) (Logger, error) {
	outputFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	return &fileLogger{
		&logger{
			level: level,
			out:   []io.WriteCloser{outputFile},
		},
		outputFile,
	}, nil
}

func (fl *fileLogger) Rotate(path string) error {
	if err := fl.Close(); err != nil {
		return err
	}

	newFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error rotating file: %w", err)
	}
	fl.file = newFile
	fl.logger.out = []io.WriteCloser{fl.file}
	return nil
}

func (fl *fileLogger) Close() error {
	return fl.file.Close()
}
