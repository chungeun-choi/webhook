package endpoint_construct

import "net/http"

type HandlerFuncManager interface {
	CreateFunc(data interface{}) (func(http.ResponseWriter, *http.Request), error)
}
