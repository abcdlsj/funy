package development

import "net/http"

type Function struct {
	Handler func(http.ResponseWriter, *http.Request)
}
