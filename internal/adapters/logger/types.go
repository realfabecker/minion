package logger

import (
	"os"
	"strings"
)

type logLevel uint8

const (
	infoLevel logLevel = iota
	errorLevel
	debugLevel
)

func (l logLevel) String() string {
	switch l {
	case debugLevel:
		return "DBG"
	case infoLevel:
		return "INF"
	case errorLevel:
		return "ERR"
	}

	return "UNKNOWN"
}

func getLogLevel() logLevel {
	switch strings.ToLower(os.Getenv("LOG")) {
	case "debug":
		return debugLevel
	default:
		return errorLevel
	}
}
