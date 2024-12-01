package logger

import (
	"fmt"
	"github.com/realfabecker/kevin/internal/core/ports"
	"log"
	"os"
)

type ConsoleLogger struct {
	level logLevel
	label string
	*log.Logger
}

func NewConsoleLogger() ports.Logger {
	logger := log.New(os.Stdout, "", 0)
	return ConsoleLogger{
		level:  getLogLevel(),
		label:  "kevin",
		Logger: logger,
	}
}

func (l ConsoleLogger) Info(message string) {
	l.log(infoLevel, message)
}

func (l ConsoleLogger) Error(message string) {
	l.log(errorLevel, message)
}

func (l ConsoleLogger) Errorf(format string, a ...interface{}) {
	l.log(errorLevel, fmt.Sprintf(format, a...))
}

func (l ConsoleLogger) Infof(format string, a ...interface{}) {
	l.log(infoLevel, fmt.Sprintf(format, a...))
}

func (l ConsoleLogger) Debug(message string) {
	l.log(debugLevel, message)
}

func (l ConsoleLogger) Debugf(format string, a ...interface{}) {
	l.log(debugLevel, fmt.Sprintf(format, a...))
}

func (l ConsoleLogger) Warn(message string) {
	l.log(warnLevel, message)
}

func (l ConsoleLogger) Fataln(v ...any) {
	l.Fatalln(v)
}

func (l ConsoleLogger) log(level logLevel, message string) {
	if l.level >= level {
		l.Logger.Printf(
			"level=\"%s\" message=\"%s\"\n",
			level,
			message,
		)
	}
}
