package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-logfmt/logfmt"
)

// Time format for logs
var TimeFormat = "2006-01-02 15:04:05"

var (
	defaultLogger = &logger{}

	Infof   = defaultLogger.Infof
	Errorf  = defaultLogger.Errorf
	Context = defaultLogger.Context
)

type Logger interface {
	Context(string) Logger
	Infof(string, ...interface{})
	Errorf(string, ...interface{})

	LogFmt() LogFmt
}

type LogFmt interface {
	Infof(...interface{})
	Errorf(...interface{})
}

var _ Logger = &logger{}
var _ LogFmt = &loggerLogFmt{}

func NewLogger() Logger {
	return &logger{
		logfmt: &loggerLogFmt{},
	}
}

func FromContext(ctx context.Context) Logger {
	// Panics on purpose, since this shouldn't happen
	return ctx.Value("log").(Logger)
}

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, "log", logger)
}

type loggerLogFmt struct {
	context string
}

func (l *loggerLogFmt) Infof(keyValues ...interface{}) {
	encoder := logfmt.NewEncoder(os.Stdout)

	encoder.EncodeKeyvals("time", time.Now().UTC().Format(TimeFormat), "level", "info", "context", l.context)
	encoder.EncodeKeyvals(keyValues...)
	encoder.EndRecord()
}

func (l *loggerLogFmt) Errorf(keyValues ...interface{}) {
	encoder := logfmt.NewEncoder(os.Stdout)

	encoder.EncodeKeyvals("time", time.Now().UTC().Format(TimeFormat), "level", "error", "context", l.context)
	encoder.EncodeKeyvals(keyValues...)
	encoder.EndRecord()
}

type logger struct {
	context string
	logfmt  *loggerLogFmt
}

func (l *logger) Context(context string) Logger {
	if l.context != "" {
		context = l.context + "|" + context
	}

	return &logger{
		context: context,
		logfmt: &loggerLogFmt{
			context: context,
		},
	}
}

func (l *logger) LogFmt() LogFmt {
	return l.logfmt
}

func (l *logger) Infof(format string, args ...interface{}) {
	args2 := append([]interface{}{time.Now().UTC().Format(TimeFormat), l.context}, args...)
	fmt.Printf("[%s] INFO [%s] "+format+"\n", args2...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	args2 := append([]interface{}{time.Now().UTC().Format(TimeFormat), l.context}, args...)
	fmt.Printf("[%s] WARN [%s] "+format+"\n", args2...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	args2 := append([]interface{}{time.Now().UTC().Format(TimeFormat), l.context}, args...)
	fmt.Printf("[%s] ERROR [%s] "+format+"\n", args2...)
}
