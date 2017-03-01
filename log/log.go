package log

import (
	"github.com/gwtony/log4go"
	"github.com/gwtony/gapi/variable"
)

// Log function
type Log interface {
	Debug(arg0 interface{}, args ...interface{})
	Info(arg0 interface{}, args ...interface{})
	Warn(arg0 interface{}, args ...interface{})
	Error(arg0 interface{}, args ...interface{})
}

// Logger logs log
type Logger struct {
	l log4go.Logger
}

// GetLogger generates logger from file
func GetLogger(path, level string, line int) *Logger {
	var log Logger

	if path == "" {
		path = variable.DEFAULT_LOG_PATH
	}

	lv := log4go.ERROR
	switch level {
	case "debug":
		lv = log4go.DEBUG
	case "info":
		lv = log4go.INFO
	case "warn":
		lv = log4go.WARNING
	case "error":
		lv = log4go.ERROR
	}

	l := log4go.NewDefaultLogger(lv)
	flw := log4go.NewFileLogWriter(path, true)
	if flw == nil {
		return nil
	}
	flw.SetFormat("[%D %T] [%L] %M")
	//flw.SetRotate(true)
	flw.SetRotateLines(line)
	//flw.SetRotateDaily(true)
	//l.AddFilter("stdout", lv, flw)
	l.AddFilter("log", lv, flw)

	log.l = l

	return &log
}

// Debug logs debug
func (l *Logger) Debug(arg0 interface{}, args ...interface{}) {
	l.l.Debug(arg0, args...)
}

// Info logs info
func (l *Logger) Info(arg0 interface{}, args ...interface{}) {
	l.l.Info(arg0, args...)
}

// Warn logs Warning
func (l *Logger) Warn(arg0 interface{}, args ...interface{}) {
	l.l.Warn(arg0, args...)
}

// Error logs error
func (l *Logger) Error(arg0 interface{}, args ...interface{}) {
	l.l.Error(arg0, args...)
}
