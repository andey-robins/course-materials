package logging

import (
	"log"
)

const LOG_LEVEL = 2

// Logger is a wrapper struct that allows us to log based on what level.
// 0 -> No Logging
// 1 -> Info
// 2 -> Warnings
// 3 -> Errors
type Logger struct {
	LogLevel int
}

func Log(logStr string, level int) {
	logger := &Logger{
		LogLevel: LOG_LEVEL,
	}
	switch level {
	case 1:
		logger.Info(logStr)
	case 2:
		logger.Warn(logStr)
	case 3:
		logger.Err(logStr)
	}
}

func (l *Logger) Info(logStr string) {
	if l.LogLevel >= 1 {
		log.Println(logStr)
	}
}

func (l *Logger) Warn(logStr string) {
	if l.LogLevel >= 2 {
		log.Println(logStr)
	}
}

func (l *Logger) Err(logStr string) {
	if l.LogLevel >= 3 {
		log.Println(logStr)
	}
}
