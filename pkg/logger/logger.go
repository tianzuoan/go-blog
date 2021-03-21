package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

type Level int8

const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_PANIC
)

func (l Level) String() string {
	switch l {
	case LEVEL_DEBUG:
		return "debug"
	case LEVEL_INFO:
		return "info"
	case LEVEL_WARN:
		return "warn"
	case LEVEL_ERROR:
		return "error"
	case LEVEL_FATAL:
		return "fatal"
	case LEVEL_PANIC:
		return "panic"
	default:
		return ""
	}
}

type Fields map[string]interface{}
type Logger struct {
	newLogger *log.Logger
	ctx       context.Context
	level     Level
	fields    Fields
	callers   []string
}

func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	logger := log.New(w, prefix, flag)
	return &Logger{newLogger: logger}
}

func (l *Logger) Debug(v ...interface{}) {
	l.WithLevel(LEVEL_DEBUG).Output(fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.WithLevel(LEVEL_DEBUG).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.WithLevel(LEVEL_INFO).Output(fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.WithLevel(LEVEL_INFO).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.WithLevel(LEVEL_WARN).Output(fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.WithLevel(LEVEL_WARN).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.WithLevel(LEVEL_ERROR).Output(fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.WithLevel(LEVEL_ERROR).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.WithLevel(LEVEL_FATAL).Output(fmt.Sprint(v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.WithLevel(LEVEL_FATAL).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(v ...interface{}) {
	l.WithLevel(LEVEL_PANIC).Output(fmt.Sprint(v...))
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.WithLevel(LEVEL_PANIC).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) clone() *Logger {
	ll := *l
	return &ll
}

// 设置日志等级
func (l *Logger) WithLevel(level Level) *Logger {
	ll := l.clone()
	ll.level = level
	return ll
}

// 设置日志公共字段
func (l *Logger) WithFields(fields Fields) *Logger {
	ll := l.clone()
	if ll.fields == nil {
		ll.fields = make(Fields)
	}
	for key, value := range fields {
		ll.fields[key] = value
	}
	return ll
}

// 设置上下文属性
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

// 设置当前某一层调用栈信息（程序计数器，文件信息和行号）
func (l *Logger) WithCaller(skip int) *Logger {
	ll := l.clone()

	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		ll.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())}
	}

	return ll
}

// 设置当前的整个调用栈信息
func (l *Logger) WithCallersFrame() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	var callers []string
	pc := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pc)
	frames := runtime.CallersFrames(pc[:depth])
	frame, more := frames.Next()
	for ; more; frame, more = frames.Next() {
		callers = append(callers, fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	ll := l.clone()
	ll.callers = callers
	return ll
}

// json格式数据
func (l *Logger) JsonFormat(message string) Fields {
	data := make(Fields, len(l.fields)+4)
	data["level"] = l.level.String()
	data["time"] = time.Now().Local().Nanosecond()
	data["message"] = message
	data["callers"] = l.callers
	if len(l.fields) > 0 {
		for key, value := range l.fields {
			if _, ok := data[key]; !ok {
				data[key] = value
			}

		}
	}
	return data
}

//真正的输出日志
func (l *Logger) Output(message string) {
	bytes, _ := json.Marshal(l.JsonFormat(message))
	content := string(bytes)
	switch l.level {
	case LEVEL_DEBUG:
		l.newLogger.Print(content)
	case LEVEL_INFO:
		l.newLogger.Print(content)
	case LEVEL_WARN:
		l.newLogger.Print(content)
	case LEVEL_FATAL:
		l.newLogger.Fatal(content)
	case LEVEL_PANIC:
		l.newLogger.Panic(content)
	default:
		l.newLogger.Print(content)
	}
}
