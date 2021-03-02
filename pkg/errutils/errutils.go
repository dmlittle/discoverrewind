package errutils

import (
	"net"
	"os"
	"syscall"

	"github.com/pkg/errors"
)

type unwrapper interface {
	Unwrap() error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Unwrap takes in an error, goes through the error cause chain, and returns the
// deepest error that has a stacktrace. This is needed because by default, every
// subsequent call to errors.Wrap overwrites the previous stacktrace. So in the
// end, we don't know where the error originated from. By going to the deepest
// one, we can find exactly where it started.
func Unwrap(err error) error {
	for {
		var u unwrapper
		if !errors.As(err, &u) {
			break
		}
		unwrapped := u.Unwrap()
		var st stackTracer
		if !errors.As(unwrapped, &st) {
			break
		}
		err = unwrapped
	}
	return err
}

// IsIgnorableErr returns true if the provided error is a EPIPE error.
func IsIgnorableErr(err error) bool {
	var serr *os.SyscallError
	if errors.As(err, &serr) {
		return serr.Err.Error() == syscall.EPIPE.Error() || serr.Err.Error() == syscall.ECONNRESET.Error()
	}

	var nerr net.Error
	if errors.As(err, &nerr) && nerr.Timeout() {
		return true
	}

	return false
}
