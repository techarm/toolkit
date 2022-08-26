package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-stack/stack"
)

const timeKey = "t"
const lvlKey = "lvl"
const msgKey = "msg"
const errorKey = "LOG15_ERROR"

// LEVEL is a type for predefined log levels.
type LEVEL int

// List of predefined log Levels
const (
	LevelTrace LEVEL = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
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

	// Log a message at the given level with context key/value pairs
	Trace(msg string, ctx ...interface{})
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
}

type logger struct {
	ctx     []interface{}
	handler *swapHandler
}

func (l *logger) write(msg string, lvl LEVEL, ctx []interface{}) {
	l.handler.Log(&Record{
		Time:    time.Now(),
		Level:   lvl,
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
	child := &logger{newContext(l.ctx, ctx), new(swapHandler)}
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

func (l *logger) Trace(msg string, ctx ...interface{}) {
	l.write(msg, LevelTrace, ctx)
}

func (l *logger) Debug(msg string, ctx ...interface{}) {
	l.write(msg, LevelDebug, ctx)
}

func (l *logger) Info(msg string, ctx ...interface{}) {
	l.write(msg, LevelInfo, ctx)
}

func (l *logger) Warn(msg string, ctx ...interface{}) {
	l.write(msg, LevelWarning, ctx)
}

func (l *logger) Error(msg string, ctx ...interface{}) {
	l.write(msg, LevelError, ctx)
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

	// ctx needs to be even because it's a series of key/value pairs
	// no one wants to check for errors on logging functions,
	// so instead of erroring on bad input, we'll just make sure
	// that things are the right length and users can fix bugs
	// when they see the output looks wrong
	if len(ctx)%2 != 0 {
		ctx = append(ctx, nil, errorKey, "Normalized odd number of arguments by adding nil")
	}

	return ctx
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
