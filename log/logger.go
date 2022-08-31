package log

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/go-stack/stack"
)

const timeKey = "t"
const lvlKey = "lvl"
const msgKey = "msg"
const errorKey = "LOG_ERROR"

// LEVEL is a type for predefined log levels.
type LEVEL int

type Fields map[string]any

// List of predefined log Levels
const (
	LevelTrace LEVEL = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

// Returns the name of a LEVEL
func (l LEVEL) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		panic("bad level")
	}
}

// LevelFromString returns the appropriate LEVEL from a string name.
// Useful for parsing command line args and configuration files.
func LevelFromString(levelString string) (LEVEL, error) {
	switch levelString {
	case "trace":
		return LevelTrace, nil
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn", "warning":
		return LevelWarning, nil
	case "error":
		return LevelError, nil
	case "fatal":
		return LevelFatal, nil
	default:
		// try to catch e.g. "INFO", "WARN" without slowing down the fast path
		lower := strings.ToLower(levelString)
		if lower != levelString {
			return LevelFromString(lower)
		}
		return LevelDebug, fmt.Errorf("log15: unknown level: %v", levelString)
	}
}

// A Record is what a Logger asks its handler to write
type Record struct {
	Time     time.Time
	Level    LEVEL
	Message  string
	Context  []interface{}
	Call     stack.Call
	KeyNames RecordKeyNames
}

// RecordKeyNames are the predefined names of the log props used by the Logger interface.
type RecordKeyNames struct {
	Time    string
	Message string
	Level   string
}

// A Logger writes key/value pairs to a Handler
type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(ctx ...interface{}) Logger

	// GetHandler gets the handler associated with the logger.
	GetHandler() Handler

	// SetHandler updates the logger to write records to the specified handler.
	SetHandler(h Handler)

	// Trace log a message at the trace level
	Trace(v ...any)

	// Debug log a message at the debug level
	Debug(v ...any)

	// Info log a message at the infomation level
	Info(v ...any)

	// Warn log a message at the warning level
	Warn(v ...any)

	// Error log a message at the error level
	Error(v ...any)

	// Fatal log a message at the fatal level
	Fatal(v ...any)

	// Tracef log a message at the trace level and arguments are handled in the manner of fmt.Printf.
	Tracef(format string, v ...any)

	// Debugf log a message at the debug level and arguments are handled in the manner of fmt.Printf.
	Debugf(format string, v ...any)

	// Infof log a message at the infomation level and arguments are handled in the manner of fmt.Printf.
	Infof(format string, v ...any)

	// Warnf log a message at the warning level and arguments are handled in the manner of fmt.Printf.
	Warnf(format string, v ...any)

	// Errorf log a message at the error level and arguments are handled in the manner of fmt.Printf.
	Errorf(format string, v ...any)

	// Fatalf log a message at the fatal level and arguments are handled in the manner of fmt.Printf.
	Fatalf(format string, v ...any)

	// Log a message at the given level with context key/value pairs
	WithField(key string, value any) Logger

	// Log a message at the given level with context key/value pairs
	WithFields(fileds Fields) Logger
}

type logger struct {
	ctx     []interface{}
	handler *swapHandler
	// fields    Fields
	// fieldPool sync.Pool
	fields atomic.Value
}

func (l *logger) write(level LEVEL, msg string) {
	fields := l.newFields()

	ctx := make([]any, 0, fields.Len())
	for _, key := range fields.Keys() {
		value, _ := fields.Get(key)
		ctx = append(ctx, key, value)
	}
	defer l.releaseFields()

	l.handler.Log(&Record{
		Time:    time.Now(),
		Level:   level,
		Message: msg,
		Context: newContext(l.ctx, ctx),
		Call:    stack.Caller(2),
		KeyNames: RecordKeyNames{
			Time:    timeKey,
			Message: msgKey,
			Level:   lvlKey,
		},
	})
}

func (l *logger) New(ctx ...interface{}) Logger {
	child := &logger{
		ctx:     newContext(l.ctx, ctx),
		handler: new(swapHandler),
		// fields:  make(map[string]any),
	}

	child.SetHandler(l.handler)
	return child
}

func newContext(prefix []interface{}, suffix []interface{}) []interface{} {
	normalizedSuffix := normalize(suffix)
	newCtx := make([]interface{}, len(prefix)+len(normalizedSuffix))
	n := copy(newCtx, prefix)
	copy(newCtx[n:], normalizedSuffix)
	return newCtx
}

func (l *logger) Trace(v ...any) {
	l.write(LevelTrace, fmt.Sprint(v...))
}

func (l *logger) Debug(v ...any) {
	l.write(LevelDebug, fmt.Sprint(v...))
}

func (l *logger) Info(v ...any) {
	l.write(LevelInfo, fmt.Sprint(v...))
}

func (l *logger) Warn(v ...any) {
	l.write(LevelWarning, fmt.Sprint(v...))
}

func (l *logger) Error(v ...any) {
	l.write(LevelError, fmt.Sprint(v...))
}

func (l *logger) Fatal(v ...any) {
	l.write(LevelFatal, fmt.Sprint(v...))
	os.Exit(1)
}

func (l *logger) Tracef(format string, args ...any) {
	l.write(LevelTrace, fmt.Sprintf(format, args...))
}

func (l *logger) Debugf(format string, args ...any) {
	l.write(LevelDebug, fmt.Sprintf(format, args...))
}

func (l *logger) Infof(format string, args ...any) {
	l.write(LevelInfo, fmt.Sprintf(format, args...))
}

func (l *logger) Warnf(format string, args ...any) {
	l.write(LevelWarning, fmt.Sprintf(format, args...))
}

func (l *logger) Errorf(format string, args ...any) {
	l.write(LevelError, fmt.Sprintf(format, args...))
}

func (l *logger) Fatalf(format string, args ...any) {
	l.write(LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *logger) WithField(key string, value any) Logger {
	fields := l.newFields()
	fields.Set(key, value)
	l.fields.Store(fields)
	return l
}

func (l *logger) WithFields(fields Fields) Logger {
	newFields := l.newFields()
	for key, value := range fields {
		newFields.Set(key, value)
	}
	l.fields.Store(newFields)
	return l
}

func (l *logger) GetHandler() Handler {
	return l.handler.Get()
}

func (l *logger) SetHandler(h Handler) {
	l.handler.Swap(h)
}

func normalize(ctx []interface{}) []interface{} {
	// if the caller passed a Context object, then expand it
	if len(ctx) == 1 {
		if ctxMap, ok := ctx[0].(Ctx); ok {
			ctx = ctxMap.toArray()
		}
	}
	return ctx
}

func (l *logger) releaseFields() {
	fields := l.fields.Load()
	if fields != nil {
		fields = orderedmap.NewOrderedMap[string, any]()
		l.fields.Store(fields)
	}
}

func (l *logger) newFields() *orderedmap.OrderedMap[string, any] {
	fields := l.fields.Load()
	if fields != nil {
		return fields.(*orderedmap.OrderedMap[string, any])
	}
	return orderedmap.NewOrderedMap[string, any]()
}

// Lazy allows you to defer calculation of a logged value that is expensive
// to compute until it is certain that it must be evaluated with the given filters.
//
// Lazy may also be used in conjunction with a Logger's New() function
// to generate a child logger which always reports the current value of changing
// state.
//
// You may wrap any function which takes no arguments to Lazy. It may return any
// number of values of any type.
type Lazy struct {
	Fn interface{}
}

// Ctx is a map of key/value pairs to pass as context to a log function
// Use this only if you really need greater safety around the arguments you pass
// to the logging functions.
type Ctx map[string]interface{}

func (c Ctx) toArray() []interface{} {
	arr := make([]interface{}, len(c)*2)

	i := 0
	for k, v := range c {
		arr[i] = k
		arr[i+1] = v
		i += 2
	}

	return arr
}
