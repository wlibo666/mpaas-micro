// Package log provides a struct logger with go-kit log
package log

import (
	"errors"
	"os"
	"strings"
	"time"

	kitlog "github.com/go-kit/kit/log"
	kitloglevel "github.com/go-kit/kit/log/level"
)

type LEVEL byte

const (
	DEBUG LEVEL = 1 << iota
	INFO
	WARN
	ERROR
)

var (
	loggers = make(map[string]kitlog.Logger)

	ErrExist    = errors.New("logger exist")
	ErrNotExist = errors.New("logger not exist")
)

func NewDefaultLogger(filename string, level ...LEVEL) error {
	err := NewLogger("default", filename, level...)
	if err != nil {
		return err
	}
	return nil
}

func NewLogger(loggerName, filename string, level ...LEVEL) error {
	_, ok := loggers[loggerName]
	if ok {
		return ErrExist
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	logger := kitlog.NewLogfmtLogger(f)

	if len(level) > 0 {
		switch level[0] {
		case DEBUG:
			logger = kitloglevel.NewFilter(logger, kitloglevel.AllowDebug())
		case INFO:
			logger = kitloglevel.NewFilter(logger, kitloglevel.AllowInfo())
		case WARN:
			logger = kitloglevel.NewFilter(logger, kitloglevel.AllowWarn())
		case ERROR:
			logger = kitloglevel.NewFilter(logger, kitloglevel.AllowError())
		default:
			logger = kitloglevel.NewFilter(logger, kitloglevel.AllowDebug())
		}
	}
	logger = kitlog.With(logger, "time", kitlog.DefaultTimestamp)
	loggers[loggerName] = logger
	go moniLoggerFile(loggerName, filename, level...)
	return nil
}

func moniLoggerFile(loggerName, filename string, level ...LEVEL) {
	for {
		_, err := os.Stat(filename)
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				reAddLogger(loggerName, filename, level...)
			}
		}
		time.Sleep(60 * time.Second)
	}
}

func reAddLogger(loggerName, filename string, level ...LEVEL) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	logger := kitlog.NewLogfmtLogger(f)

	if len(level) > 0 {
		switch level[0] {
		case DEBUG:
			logger = kitloglevel.NewFilter(logger, kitloglevel.AllowDebug())
		case INFO:
			logger = kitloglevel.NewFilter(logger, kitloglevel.AllowInfo())
		case WARN:
			logger = kitloglevel.NewFilter(logger, kitloglevel.AllowWarn())
		case ERROR:
			logger = kitloglevel.NewFilter(logger, kitloglevel.AllowError())
		default:
			logger = kitloglevel.NewFilter(logger, kitloglevel.AllowDebug())
		}
	}
	logger = kitlog.With(logger, "time", kitlog.DefaultTimestamp)
	loggers[loggerName] = logger
	return nil
}

func Debug(keyvals ...interface{}) error {
	logger, ok := loggers["default"]
	if !ok {
		return ErrNotExist
	}
	return kitloglevel.Debug(logger).Log(keyvals...)
}

func Info(keyvals ...interface{}) error {
	logger, ok := loggers["default"]
	if !ok {
		return ErrNotExist
	}
	return kitloglevel.Info(logger).Log(keyvals...)
}

func Warn(keyvals ...interface{}) error {
	logger, ok := loggers["default"]
	if !ok {
		return ErrNotExist
	}
	return kitloglevel.Warn(logger).Log(keyvals...)
}

func Error(keyvals ...interface{}) error {
	logger, ok := loggers["default"]
	if !ok {
		return ErrNotExist
	}
	return kitloglevel.Error(logger).Log(keyvals...)
}

func DebugWithLogger(loggerName string, keyvals ...interface{}) error {
	logger, ok := loggers[loggerName]
	if !ok {
		return ErrNotExist
	}
	return kitloglevel.Debug(logger).Log(keyvals...)
}

func InfoWithLogger(loggerName string, keyvals ...interface{}) error {
	logger, ok := loggers[loggerName]
	if !ok {
		return ErrNotExist
	}
	return kitloglevel.Info(logger).Log(keyvals...)
}

func WarnWithLogger(loggerName string, keyvals ...interface{}) error {
	logger, ok := loggers[loggerName]
	if !ok {
		return ErrNotExist
	}
	return kitloglevel.Warn(logger).Log(keyvals...)
}

func ErrorWithLogger(loggerName string, keyvals ...interface{}) error {
	logger, ok := loggers[loggerName]
	if !ok {
		return ErrNotExist
	}
	return kitloglevel.Error(logger).Log(keyvals...)
}

// ExampleLogger is useless,just for showing example code by godoc
type ExampleLogger struct{}
