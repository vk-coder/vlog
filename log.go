// Package vlog implements a simple logging package on top of existing Go log package.
// vlog package adds logging levels and allows changing the logging levels. Otherwise
// vlog package behaves exactly as the Go log package.
// vlog does not create any hirarchy of the loggers.
// Root logger is always created with name "".
package vlog

import (
	"fmt"
	"io"
	"log"
	"sync"
)

type Level int

// Logging levels
const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	CRITICAL
)

type Logger struct {
	logger *log.Logger
	name   string
	level  Level
}

var defaultFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
var rootLogger *Logger

var mu sync.Mutex
var loggers map[string]*Logger

func init() {
	loggers = make(map[string]*Logger)

	rootLogger = &Logger{
		name:   "",
		level:  TRACE,
		logger: log.Default(),
	}

	rootLogger.logger.SetFlags(defaultFlags)
	loggers[""] = rootLogger
}

func newLogger(name string, out io.Writer, flag int) *Logger {
	if len(name) == 0 {
		return loggers[""]
	}

	rootLogger = &Logger{
		name:   name,
		level:  TRACE,
		logger: log.New(out, name+" ", flag),
	}
	loggers[name] = rootLogger

	return rootLogger
}

// GetLogger creates a new Logger with specified name.
// The new logger inherits properties from the parent logger.
// If parent logger is nil then new logger inherits properties from
// root logger.
func GetLogger(name string, parent *Logger) *Logger {
	mu.Lock()
	defer mu.Unlock()

	if logger, found := loggers[name]; found {
		return logger
	}

	if parent != nil {
		return newLogger(name, parent.logger.Writer(), parent.logger.Flags())
	}

	rootLogger := loggers[""]
	return newLogger(name, rootLogger.logger.Writer(), rootLogger.logger.Flags())
}

// SetFlags proxy call to log.Logger SetFlags
func (l *Logger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

// SetOutput proxy call to log.Logger SetOutput
func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

// GetOutput returns the current output destination
func (l *Logger) GetOutput() io.Writer {
	return l.logger.Writer()
}

// SetLevel sets the logging level for the logger.
// All the levels above given level would be emitted.
func (l *Logger) SetLevel(lvl Level) {
	l.level = lvl
}

// GetLevel returns the current logging level
func (l *Logger) GetLevel() Level {
	return l.level
}

func (l *Logger) output(level string, v ...interface{}) {
	l.logger.Output(3, fmt.Sprint(append([]interface{}{level}, v...)...))
}

// Trace emits a log at trace logging level
func (l *Logger) Trace(v ...interface{}) {
	if TRACE >= l.level {
		l.output("trace ", v...)
	}
}

// Debug emits a log at debug logging level
func (l *Logger) Debug(v ...interface{}) {
	if DEBUG >= l.level {
		l.output("debug ", v...)
	}
}

// Info emits a log at info logging level
func (l *Logger) Info(v ...interface{}) {
	if INFO >= l.level {
		l.output("info ", v...)
	}
}

// Warn emits a log at warn logging level
func (l *Logger) Warn(v ...interface{}) {
	if WARN >= l.level {
		l.output("warn ", v...)
	}
}

// Error emits a log at error logging level
func (l *Logger) Error(v ...interface{}) {
	if ERROR >= l.level {
		l.output("error ", v...)
	}
}

// Critical emits a log at critical logging level
func (l *Logger) Critical(v ...interface{}) {
	if CRITICAL >= l.level {
		l.output("critical ", v...)
	}
}

// Tracef emits a log at tracef logging level with specified format string
func (l *Logger) Tracef(format string, v ...interface{}) {
	if TRACE >= l.level {
		l.output("trace ", fmt.Sprintf(format, v...))
	}
}

// Debugf emits a log at debugf logging level with specified format string
func (l *Logger) Debugf(format string, v ...interface{}) {
	if DEBUG >= l.level {
		l.output("debug ", fmt.Sprintf(format, v...))
	}
}

// Infof emits a log at infof logging level with specified format string
func (l *Logger) Infof(format string, v ...interface{}) {
	if INFO >= l.level {
		l.output("info ", fmt.Sprintf(format, v...))
	}
}

// Warnf emits a log at warnf logging level with specified format string
func (l *Logger) Warnf(format string, v ...interface{}) {
	if WARN >= l.level {
		l.output("warn ", fmt.Sprintf(format, v...))
	}
}

// Errorf emits a log at errorf logging level with specified format string
func (l *Logger) Errorf(format string, v ...interface{}) {
	if ERROR >= l.level {
		l.output("error ", fmt.Sprintf(format, v...))
	}
}

// Criticalf emits a log at criticalf logging level with specified format string
func (l *Logger) Criticalf(format string, v ...interface{}) {
	if CRITICAL >= l.level {
		l.output("critical ", fmt.Sprintf(format, v...))
	}
}
