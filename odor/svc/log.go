//
// Copyright (c) Telefonica I+D. All rights reserved.
//

package svc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"
)

type level int

const (
	debugLevel level = iota
	infoLevel
	warnLevel
	errorLevel
	fatalLevel
)

var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

func levelByName(levelName string) level {
	levelName = strings.ToUpper(levelName)
	for i, name := range levelNames {
		if name == levelName {
			return level(i)
		}
	}
	return infoLevel
}

// Logger type.
type Logger struct {
	out      io.Writer
	logLevel level
	context  interface{}
	mutex    sync.Mutex
}

// NewLogger to create a Logger.
func NewLogger() *Logger {
	return &Logger{
		out:      os.Stdout,
		logLevel: infoLevel,
	}
}

// SetLogContext to set a global context.
func (l *Logger) SetLogContext(context interface{}) {
	l.context = context
}

// GetLogContext to get the global context.
func (l *Logger) GetLogContext() interface{} {
	return l.context
}

// SetLevel to set the log level.
func (l *Logger) SetLevel(levelName string) {
	l.logLevel = levelByName(levelName)
}

// GetLevel to return the log level.
func (l *Logger) GetLevel() string {
	return levelNames[l.logLevel]
}

func (l *Logger) log(logLevel level, context interface{}, message string, args ...interface{}) {
	if logLevel >= l.logLevel {
		var buf bytes.Buffer
		writeDoc(&buf, time.Now(), levelNames[logLevel], l.context, context, fmt.Sprintf(message, args...))
		bytes := buf.Bytes()
		l.mutex.Lock()
		defer l.mutex.Unlock()
		l.out.Write(bytes)
	}
}

func writeDoc(buf *bytes.Buffer, time time.Time, level string, context, customContext interface{}, message string) {
	buf.WriteByte('{')
	writeField(buf, "time", time)
	buf.WriteByte(',')
	writeField(buf, "lvl", level)
	buf.WriteByte(',')
	if length := writeFields(buf, context); length > 0 {
		buf.WriteByte(',')
	}
	if length := writeFields(buf, customContext); length > 0 {
		buf.WriteByte(',')
	}
	writeField(buf, "msg", message)
	buf.WriteByte('}')
	buf.WriteByte('\n')
}

func writeField(buf *bytes.Buffer, key string, value interface{}) {
	buf.WriteByte('"')
	buf.WriteString(key)
	buf.WriteByte('"')
	buf.WriteByte(':')
	if jsonValue, err := json.Marshal(value); err == nil {
		buf.Write(jsonValue)
	}
}

func writeFields(buf *bytes.Buffer, fields interface{}) int {
	if fields == nil {
		return 0
	}
	if jsonFields, err := json.Marshal(fields); err == nil {
		if len(jsonFields) > 2 && jsonFields[0] == '{' && jsonFields[len(jsonFields)-1] == '}' {
			if _, err := buf.Write(jsonFields[1 : len(jsonFields)-1]); err == nil {
				return len(jsonFields) - 2
			}
		}
	}
	return 0
}

// Debug to log a message at debug level
func (l *Logger) Debug(message string, args ...interface{}) {
	l.log(debugLevel, nil, message, args...)
}

// DebugC to log a message at debug level with custom context
func (l *Logger) DebugC(context interface{}, message string, args ...interface{}) {
	l.log(debugLevel, context, message, args...)
}

// Info to log a message at info level
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(infoLevel, nil, message, args...)
}

// InfoC to log a message at info level
func (l *Logger) InfoC(context interface{}, message string, args ...interface{}) {
	l.log(infoLevel, context, message, args...)
}

// Warn to log a message at warn level
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(warnLevel, nil, message, args...)
}

// WarnC to log a message at warn level
func (l *Logger) WarnC(context interface{}, message string, args ...interface{}) {
	l.log(warnLevel, context, message, args...)
}

// Error to log a message at error level
func (l *Logger) Error(message string, args ...interface{}) {
	l.log(errorLevel, nil, message, args...)
}

// ErrorC to log a message at error level
func (l *Logger) ErrorC(context interface{}, message string, args ...interface{}) {
	l.log(errorLevel, context, message, args...)
}

// Fatal to log a message at fatal level
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.log(fatalLevel, nil, message, args...)
}

// FatalC to log a message at fatal level
func (l *Logger) FatalC(context interface{}, message string, args ...interface{}) {
	l.log(fatalLevel, context, message, args...)
}

// DebugResponse to dump the response at debug level.
func (l *Logger) DebugResponse(message string, r *http.Response) {
	l.DebugResponseC(nil, message, r)
}

// DebugResponseC to dump the response at debug level.
func (l *Logger) DebugResponseC(context interface{}, message string, r *http.Response) {
	if r != nil && l.logLevel <= debugLevel {
		if dump, err := httputil.DumpResponse(r, true); err == nil {
			l.DebugC(context, "%s. %s", message, dump)
		}
	}
}

// DebugRequest to dump the request at debug level.
func (l *Logger) DebugRequest(message string, r *http.Request) {
	l.DebugRequestC(nil, message, r)
}

// DebugRequestC to dump the request at debug level.
func (l *Logger) DebugRequestC(context interface{}, message string, r *http.Request) {
	if r != nil && l.logLevel <= debugLevel {
		if dump, err := httputil.DumpRequest(r, true); err == nil {
			l.DebugC(context, "%s. %s", message, dump)
		}
	}
}

// DebugRequestOut to dump the output request at debug level.
func (l *Logger) DebugRequestOut(message string, r *http.Request) {
	l.DebugRequestOutC(nil, message, r)
}

// DebugRequestOutC to dump the output request at debug level.
func (l *Logger) DebugRequestOutC(context interface{}, message string, r *http.Request) {
	if r != nil && l.logLevel <= debugLevel {
		if dump, err := httputil.DumpRequestOut(r, true); err == nil {
			l.DebugC(context, "%s. %s", message, dump)
		}
	}
}

// Bridge to std log
type writer struct {
	l *Logger
}

func (w *writer) Write(p []byte) (int, error) {
	w.l.Error(string(p))
	return len(p), nil
}

// NewStdLogger returns a standard logger struct but using our custom logger.
func NewStdLogger(l *Logger) *log.Logger {
	sl := log.New(&writer{l: l}, "", 0)
	return sl
}

// NewStdLoggerC returns a standard logger struct but using our custom logger with a specific context.
func NewStdLoggerC(context interface{}) *log.Logger {
	logger := NewLogger()
	logger.SetLogContext(context)
	return NewStdLogger(logger)
}
