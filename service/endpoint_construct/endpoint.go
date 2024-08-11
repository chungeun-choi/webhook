package endpoint_construct

import (
	"github.com/chungeun-choi/webhook/common/errors"
	"net/http"
	"sync"
)

// EndpointManager is a struct that manages the endpoints
type EndpointManager struct {
	Server  *http.Server
	Mux     *http.ServeMux
	muxLock sync.RWMutex
}

// NewEndpointManager creates a new endpoint manager
func NewEndpointManager(server *http.Server, endpointGroup string) *EndpointManager {
	if server.Handler == nil {
		mux := http.NewServeMux()
		server.Handler = mux
	}

	return &EndpointManager{
		Server: server,
		Mux:    server.Handler.(*http.ServeMux),
	}
}

func (w *EndpointManager) AddEndpoint(endpoint string, constructor HandlerFuncManager, data interface{}) error {
	return w.addEndpoint(endpoint, constructor, data)
}

// addEndpoint adds an endpoint to the server
func (w *EndpointManager) addEndpoint(endpoint string, constructor HandlerFuncManager, data interface{}) error {
	if constructor == nil {
		return errors.ErrNotFound
	}

	var (
		handlerFunc func(http.ResponseWriter, *http.Request)
		err         error
	)

	// Create the handler function
	if handlerFunc, err = constructor.CreateFunc(data); err != nil {
		return err
	}

	// Add the handler function to the endpoint
	w.muxLock.Lock()
	defer w.muxLock.Unlock()
	w.Mux.HandleFunc(endpoint, handlerFunc)

	return nil
}
