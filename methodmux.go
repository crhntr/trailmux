package trailmux

import "net/http"

type MethodMux struct {
	GET, POST, DELETE, PUT, PATCH, HEAD, CONNECT, OPTIONS, TRACE http.Handler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler
}

func (methMux MethodMux) serveHttp(h http.Handler, w http.ResponseWriter, r *http.Request) {
	if h == nil {
		if methMux.MethodNotAllowed != nil {
			methMux.MethodNotAllowed.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}
	h.ServeHTTP(w, r)
}

func (methMux MethodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		methMux.serveHttp(methMux.GET, w, r)
	case "POST":
		methMux.serveHttp(methMux.POST, w, r)
	case "DELETE":
		methMux.serveHttp(methMux.DELETE, w, r)
	case "PUT":
		methMux.serveHttp(methMux.PUT, w, r)
	case "PATCH":
		methMux.serveHttp(methMux.PATCH, w, r)
	case "HEAD":
		methMux.serveHttp(methMux.HEAD, w, r)
	case "CONNECT":
		methMux.serveHttp(methMux.CONNECT, w, r)
	case "OPTIONS":
		methMux.serveHttp(methMux.OPTIONS, w, r)
	case "TRACE":
		methMux.serveHttp(methMux.TRACE, w, r)
	default:
		http.Error(w,
			http.StatusText(http.StatusBadRequest)+" Method "+r.Method+" Unknown",
			http.StatusBadRequest,
		)
	}
}
