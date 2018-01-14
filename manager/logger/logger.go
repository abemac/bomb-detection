package logger

import (
	"log"
	"os"

	"github.com/abemac/bomb-detection/manager/constants"
)

type Logger struct {
	logger *log.Logger
}

const (
	OFF     = 0
	FATAL   = 1
	ERROR   = 2
	WARN    = 3
	INFO    = 4
	DEBUG   = 5
	VERBOSE = 6
)

func NewLogger(name string) *Logger {
	l := new(Logger)
	l.logger = log.New(os.Stdout, name+" ", log.Ltime|log.Lmicroseconds)
	return l
}

func (l *Logger) V(line ...interface{}) {
	if constants.LOG_LEVEL >= VERBOSE {
		l.logger.Println("VERBOSE: ", line)
	}
}
func (l *Logger) D(line ...interface{}) {
	if constants.LOG_LEVEL >= DEBUG {
		l.logger.Println("DEBUG: ", line)
	}
}
func (l *Logger) I(line ...interface{}) {
	if constants.LOG_LEVEL >= INFO {
		l.logger.Println("INFO: ", line)
	}
}
func (l *Logger) W(line ...interface{}) {
	if constants.LOG_LEVEL >= WARN {
		l.logger.Println("WARN: ", line)
	}
}
func (l *Logger) E(line ...interface{}) {
	if constants.LOG_LEVEL >= ERROR {
		l.logger.Println("ERROR: ", line)
	}
}
func (l *Logger) F(line ...interface{}) {
	if constants.LOG_LEVEL >= FATAL {
		l.logger.Println("FATAL: ", line)
	}
}
