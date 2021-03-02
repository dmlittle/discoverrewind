package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/dmlittle/discoverrewind/pkg/errutils"
	"github.com/dmlittle/discoverrewind/pkg/logger/writer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Data is a type alias so that it is easy to add additional data to log lines.
type Data map[string]interface{}

// Logger is a logger instance that contains necessary info needed when logging.
type Logger struct {
	zl   zerolog.Logger
	id   string
	data []Data
	err  error
	root []Data
}

const stackSize = 4 << 10 // 4KB

var (
	output io.Writer = os.Stdout
	host   string
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func init() {
	host, _ = os.Hostname() // nolint: gosec
	zerolog.TimestampFieldName = "timestamp"
	output = writer.ConsoleWriter{Out: os.Stderr}
}

// New returns a new configured Logger instance.
func New() Logger {
	return newWithLevel(os.Getenv("LOG_LEVEL"))
}

func NewWithLevel(level string) Logger {
	return newWithLevel(level)
}

func newWithLevel(level string) Logger {
	zl := zerolog.New(output).
		With().
		Timestamp().
		Str("hostname", host).
		Logger()

	switch level {
	case "debug":
		zl = zl.Level(zerolog.DebugLevel)
	case "", "info":
		zl = zl.Level(zerolog.InfoLevel)
	case "warn":
		zl = zl.Level(zerolog.WarnLevel)
	case "error":
		zl = zl.Level(zerolog.ErrorLevel)
	}

	return Logger{
		zl:   zl,
		data: []Data{},
		root: []Data{},
	}
}

func (log Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, log)
}

// ID adds an identifier to every subsequent log line.
func (log Logger) ID(id string) Logger {
	log.id = id
	return log
}

// Data adds additional data to every subsequent log line. This data is nested
// within the "data" field.
func (log Logger) Data(data Data) Logger {
	log.data = append(log.data, data)
	return log
}

// Err adds additional an error object to every subsequent log line. This is
// meant to be immediately chained to emit the log line.
func (log Logger) Err(err error) Logger {
	log.err = err
	return log
}

// Root adds additional data to every subsequent log line. This data is added to
// the root of the JSON line.
func (log Logger) Root(root Data) Logger {
	log.root = append(log.root, root)
	return log
}

// Info emits a log line with an "info" log level.
func (log Logger) Info(message string, fields ...Data) {
	log.log(log.zl.Info(), message, fields...)
}

// Error emits a log line with an "error" log level.
func (log Logger) Error(message string, fields ...Data) {
	log.log(log.zl.Error(), message, fields...)
}

// Warn emits a log line with an "warn" log level.
func (log Logger) Warn(message string, fields ...Data) {
	log.log(log.zl.Warn(), message, fields...)
}

// Debug emits a log line with an "debug" log level.
func (log Logger) Debug(message string, fields ...Data) {
	log.log(log.zl.Debug(), message, fields...)
}

// Fatal emits a log line with an "fatal" log level. It also calls os.Exit(1)
// afterwards.
func (log Logger) Fatal(message string, fields ...Data) {
	log.log(log.zl.Fatal(), message, fields...)
}

func (log Logger) log(evt *zerolog.Event, message string, fields ...Data) {
	// Merge data fields
	hasData := false
	data := zerolog.Dict()
	for _, field := range append(log.data, fields...) {
		if len(field) != 0 {
			hasData = true
			data = data.Fields(field)
		}
	}

	// Add root fields
	for _, field := range log.root {
		if len(field) != 0 {
			evt = evt.Fields(field)
		}
	}
	// Add id field
	if log.id != "" {
		evt = evt.Str("id", log.id)
	}
	// Add data field
	if hasData {
		evt = evt.Dict("data", data)
	}

	if log.err != nil {
		var stack []byte

		// This is to support pkg/errors stackTracer interface.
		err := errutils.Unwrap(log.err)
		var st stackTracer
		if errors.As(err, &st) {
			stack = []byte(fmt.Sprintf("%+v", st.StackTrace()))
		} else {
			stack = make([]byte, stackSize)
			n := runtime.Stack(stack, true)
			stack = stack[:n]
		}
		evt = evt.Str("error", log.err.Error())
		evt = evt.Bytes("stack", stack)
	}

	evt.Int64("nanoseconds", time.Now().UnixNano()).Msg(message)
}
