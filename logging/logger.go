package logging

import (
	"bytes"
	"fmt"
	"log"
)

type LogLevel int

var isInitialized bool
var cLogger chagallLogger

const (
	NO_LOGGING LogLevel = iota
	LOG_FATAL
	LOG_WARN
	LOG_ALL
)

//TODO: Add the possibility to log into a file, not only on stdout
type chagallLogger struct {
	logger   *log.Logger
	loglevel LogLevel
	buff     *bytes.Buffer
}

func InitLogger(logLevel LogLevel) {
	var buff bytes.Buffer
	goLogger := log.New(&buff, "INFO: ", log.LstdFlags|log.Llongfile)
	cLogger = chagallLogger{
		logger:   goLogger,
		loglevel: logLevel,
		buff:     &buff,
	}
}

//When the log level is too low, the log function will simply return and log nothing
func LogInfo(msg string) {

	if cLogger.loglevel < LOG_ALL {
		return
	}

	cLogger.logger.SetPrefix("[INFO]\t")
	cLogger.logger.Output(2, msg)
	fmt.Print(cLogger.buff)
	cLogger.buff.Reset()
}

func LogWarning(msg string) {

	if cLogger.loglevel < LOG_WARN {
		return
	}

	cLogger.logger.SetPrefix("[WARN]\t")
	cLogger.logger.Output(2, msg)
	fmt.Print(cLogger.buff)
	cLogger.buff.Reset()
}

func LogFatalError(msg string) {

	if cLogger.loglevel < LOG_FATAL {
		return
	}

	cLogger.logger.SetPrefix("[FATAL]\t")
	cLogger.logger.Output(2, msg)
	fmt.Print(cLogger.buff)
	cLogger.buff.Reset()
}
