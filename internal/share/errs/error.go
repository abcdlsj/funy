package errs

import "errors"

var (
	// Dynamic library loader error
	ErrDyLibLoadNoRun      = errors.New("no run function in dylib")
	ErrDyLibLoadNoShutdown = errors.New("no shutdown function in dylib")
)
