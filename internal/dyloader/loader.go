package dyloader

import (
	"plugin"

	"github.com/abcdlsj/funy/internal/share/errs"
	"github.com/charmbracelet/log"
)

type fn func() error

const RunFnName = "Run"
const ShutdownFnName = "Shutdown"

func Load(dylib string) (run, shutdown fn, err error) {
	p, err := plugin.Open(dylib)
	if err != nil {
		log.Errorf("plugin open error: %v", err)
		return
	}

	runSym, err := p.Lookup(RunFnName)
	if err != nil {
		log.Errorf("plugin lookup error: %v", err)
		return
	}

	if runFn, ok := runSym.(func() error); ok {
		run = runFn
	} else {
		err = errs.ErrDyLibLoadNoRun
		return
	}

	shutdownSym, err := p.Lookup(RunFnName)
	if err != nil {
		log.Errorf("plugin lookup error: %v", err)
		return
	}

	if shutdownFn, ok := shutdownSym.(func() error); ok {
		shutdown = shutdownFn
	} else {
		err = errs.ErrDyLibLoadNoRun
		return
	}

	return
}
