package dyloader

import (
	"net/http"
	"plugin"

	"github.com/abcdlsj/funy/pkgs/share/errs"
	"github.com/charmbracelet/log"
)

const HandlerName = "Handler"

func LoadFunction(dylib string) (handler func(http.ResponseWriter, *http.Request), err error) {
	p, err := plugin.Open(dylib)
	if err != nil {
		log.Errorf("plugin open error: %v", err)
		return
	}

	handlerSym, err := p.Lookup(HandlerName)
	if err != nil {
		log.Errorf("plugin lookup error: %v", err)
		return
	}

	if handlerFn, ok := handlerSym.(func(http.ResponseWriter, *http.Request)); ok {
		handler = handlerFn
	} else {
		err = errs.ErrDyLibLoadNoRun
		return
	}

	return
}
