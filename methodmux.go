package trailmux

import (
	"fmt"
	"net/http"
)

var methodStrings = [...]string{"GET", "POST", "DELETE", "PUT", "PATCH", "HEAD", "CONNECT", "OPTIONS", "TRACE"}

func IsMethod(methodCandidate string) bool {
	for _, method := range methodStrings {
		if method == methodCandidate {
			return true
		}
	}
	return false
}

type MethodMux struct {
	methodHandlers [len(methodStrings)]http.Handler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler
}

func (methMux MethodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for i, method := range methodStrings {
		if method == r.Method {
			handler := methMux.methodHandlers[i]
			if handler == nil {
				break
			}
			handler.ServeHTTP(w, r)
		}
	}

	if methMux.MethodNotAllowed != nil {
		methMux.MethodNotAllowed.ServeHTTP(w, r)
	}

	http.Error(w,
		http.StatusText(http.StatusMethodNotAllowed),
		http.StatusMethodNotAllowed,
	)
}

func (methMux *MethodMux) Handle(methodToHandle string, handler http.Handler) {
	if !IsMethod(methodToHandle) {
		panic("Not a valid http method")
	}
	for i, method := range methodStrings {
		if methodToHandle == method {
			if methMux.methodHandlers[i] != nil {
				panic(fmt.Sprintf("handler for %s already set", method))
			}
			methMux.methodHandlers[i] = handler
			break
		}
	}
}

func (mux *MethodMux) GET(handler http.Handler) {
	mux.Handle("GET", handler)
}
func (mux *MethodMux) POST(handler http.Handler) {
	mux.Handle("POST", handler)
}
func (mux *MethodMux) DELETE(handler http.Handler) {
	mux.Handle("DELETE", handler)
}
func (mux *MethodMux) PUT(handler http.Handler) {
	mux.Handle("PUT", handler)
}
func (mux *MethodMux) PATCH(handler http.Handler) {
	mux.Handle("PATCH", handler)
}
func (mux *MethodMux) HEAD(handler http.Handler) {
	mux.Handle("HEAD", handler)
}
func (mux *MethodMux) CONNECT(handler http.Handler) {
	mux.Handle("CONNECT", handler)
}
func (mux *MethodMux) OPTIONS(handler http.Handler) {
	mux.Handle("OPTIONS", handler)
}
func (mux *MethodMux) TRACE(handler http.Handler) {
	mux.Handle("TRACE", handler)
}
