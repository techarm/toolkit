package log

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

// Predefined handlers
var (
	root          *logger
	StdoutHandler = StreamHandler(os.Stdout, LogfmtFormat())
)

func init() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		StdoutHandler = StreamHandler(os.Stdout, TerminalFormatterDefault())
	}

	root = &logger{
		ctx:     []interface{}{},
		handler: new(swapHandler),
		// fields:  make(map[string]any),
	}
	root.SetHandler(StdoutHandler)
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func New(ctx ...interface{}) Logger {
	return root.New(ctx...)
}

// Root returns the root logger
func Root() Logger {
	return root
}

// The following functions bypass the exported logger methods (logger.Debug,
// etc.) to keep the call depth the same for all paths to logger.write so
// runtime.Caller(2) always refers to the call site in client code.

// Trace is a convenient alias for Root().Trace
func Trace(v ...any) {
	root.write(LevelTrace, fmt.Sprint(v...))
}

// Debug is a convenient alias for Root().Debug
func Debug(v ...any) {
	root.write(LevelDebug, fmt.Sprint(v...))
}

// Info is a convenient alias for Root().Info
func Info(v ...any) {
	root.write(LevelInfo, fmt.Sprint(v...))
}

// Warn is a convenient alias for Root().Warn
func Warn(v ...any) {
	root.write(LevelWarning, fmt.Sprint(v...))
}

// Error is a convenient alias for Root().Error
func Error(v ...any) {
	root.write(LevelError, fmt.Sprint(v...))
}

// Fatal is a convenient alias for Root().Fatal
func Fatal(v ...any) {
	root.write(LevelFatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Tracef is a convenient alias for Root().Debugf
func Tracef(format string, v ...any) {
	root.write(LevelTrace, fmt.Sprintf(format, v...))
}

// Debugf is a convenient alias for Root().Debugf
func Debugf(format string, v ...any) {
	root.write(LevelDebug, fmt.Sprintf(format, v...))
}

// Infof is a convenient alias for Root().Infof
func Infof(format string, v ...any) {
	root.write(LevelInfo, fmt.Sprintf(format, v...))
}

// Warnf is a convenient alias for Root().Warnf
func Warnf(format string, v ...any) {
	root.write(LevelWarning, fmt.Sprintf(format, v...))

}

// Errorf is a convenient alias for Root().Errorf
func Errorf(format string, v ...any) {
	root.write(LevelError, fmt.Sprintf(format, v...))
}

// Fatalf is a convenient alias for Root().Fatalf
func Fatalf(format string, v ...any) {
	root.write(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// WithField is a convenient alias for Root().WithField
func WithField(k string, v any) Logger {
	return root.WithField(k, v)
}

// WithFields is a convenient alias for Root().WithFields
func WithFields(fields Fields) Logger {
	return root.WithFields(fields)
}
