package flexilog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type LogLevel uint8

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger interface {
	Debug(format string, args ...any) error
	Info(format string, args ...any) error
	Warn(format string, args ...any) error
	Error(format string, args ...any) error
	Fatal(format string, args ...any)
	Close() error
}

type logger struct {
	mu    sync.Mutex
	level LogLevel
	out   []io.WriteCloser
}

func New(level LogLevel, outs ...io.WriteCloser) Logger {
	return &logger{
		level: level,
		out:   outs,
	}
}

func (l *logger) Debug(format string, args ...any) error {
	return l.log(DEBUG, "[DEBUG]:", format, args...)
}

func (l *logger) Info(format string, args ...any) error {
	return l.log(INFO, "[INFO]: ", format, args...)
}

func (l *logger) Warn(format string, args ...any) error {
	return l.log(WARN, "[WARN]: ", format, args...)
}

func (l *logger) Error(format string, args ...any) error {
	return l.log(ERROR, "[ERROR]:", format, args...)
}

func (l *logger) Fatal(format string, args ...any) {
	err := l.log(FATAL, "[FATAL]:", format, args...)
	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			panic(err)
		}
	}
	os.Exit(1)
}

func (l *logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, out := range l.out {
		err := out.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *logger) log(level LogLevel, prefix, format string, args ...any) error {
	if level < l.level {
		return fmt.Errorf("log level too low: %v", level)
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	color := color(level)
	msg := fmt.Sprintf(format, args...)
	for _, outer := range l.out {
		_, err := fmt.Fprintf(
			outer, "%s%s%s %s%s\n",
			string(color),
			time.Now().Format("2006-01-02 15:04:05 "),
			prefix,
			msg,
			RESET,
		)
		if err != nil {
			return fmt.Errorf("writing log: %s\n", err)
		}
	}

	return nil
}

type colorType string

const (
	RED    colorType = "\u001B[31m"
	GREEN  colorType = "\u001B[32m"
	YELLOW colorType = "\u001B[33m"
	PURPLE colorType = "\u001B[35m"
	CYAN   colorType = "\u001B[36m"
	RESET  colorType = "\u001B[0m"
)

func color(level LogLevel) colorType {
	switch level {
	case DEBUG:
		return CYAN
	case INFO:
		return GREEN
	case WARN:
		return YELLOW
	case ERROR:
		return RED
	case FATAL:
		return PURPLE
	default:
		return RESET
	}
}
