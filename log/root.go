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

	root = &logger{[]interface{}{}, new(swapHandler)}
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
func Trace(msg string, ctx ...interface{}) {
	root.write(msg, LevelTrace, ctx)
}

// Debug is a convenient alias for Root().Debug
func Debug(msg string, ctx ...interface{}) {
	root.write(msg, LevelDebug, ctx)
}

// Info is a convenient alias for Root().Info
func Info(msg string, ctx ...interface{}) {
	root.write(msg, LevelInfo, ctx)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, ctx ...interface{}) {
	root.write(msg, LevelWarning, ctx)
}

// Error is a convenient alias for Root().Error
func Error(msg string, ctx ...interface{}) {
	root.write(msg, LevelError, ctx)
}

// Tracef formats according to a format specifier and returns the resulting string.
func Tracef(format string, args ...any) {
	root.write(fmt.Sprintf(format, args...), LevelTrace, nil)
}

// Debugf is a convenient alias for Root().Debugf
func Debugf(format string, args ...any) {
	root.write(fmt.Sprintf(format, args...), LevelDebug, nil)
}

// Infof is a convenient alias for Root().Infof
func Infof(format string, args ...any) {
	root.write(fmt.Sprintf(format, args...), LevelInfo, nil)
}

// Warnf is a convenient alias for Root().Warnf
func Warnf(format string, args ...any) {
	root.write(fmt.Sprintf(format, args...), LevelWarning, nil)
}

// Errorf is a convenient alias for Root().Errorf
func Errorf(format string, args ...any) {
	root.write(fmt.Sprintf(format, args...), LevelError, nil)
}

// Errors is a convenient alias for Root().Errorf
func Errors(err error) {
	root.write(err.Error(), LevelError, nil)
}

// Errorv is a convenient alias for Root().Errorf
func Errorv(err error) {
	root.write(fmt.Sprintf("%v", err), LevelError, nil)
}
