package trailmux

import (
  "net/http"
  "strings"
)

type Routes map[string]http.Handler

func (routes Routes) Mux() Mux {
  return NewMux(routes)
}

type Mux struct {
  methods map[string]http.Handler
  paths map[string]http.Handler

  NoMatch http.Handler
}

func (mux Mux) NoMatchHandler(handler http.Handler) Mux {
  mux.NoMatch = handler
  return mux
}

func (mux Mux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  for path, handler := range mux.paths {
    if strings.HasPrefix(req.URL.Path, path) {
      http.StripPrefix(path, handler).ServeHTTP(res, req)
      return
    }
  }

  for method, handler := range mux.methods {
    if req.Method == method {
      handler.ServeHTTP(res, req)
      return
    }
  }

  mux.NoMatch.ServeHTTP(res, req)
}

// NewMux sorts routes by method or path
// sets NoMatch to a default handler writing HTTP status 404
// when any path is added or when no method handlers are added.
// A method not allowed handler is maped to NoMatch when all handlers
// are HTTP Method strings
func NewMux(routes map[string]http.Handler) Mux {
  var mux Mux
  mux.methods, mux.paths = sortRoutes(routes)

  if len(mux.paths) != 0 || len(mux.methods) == 0 {
    mux.NoMatch = http.HandlerFunc(defaultNotFound)
  } else {
    mux.NoMatch = http.HandlerFunc(defaultMethodNotAllowed)
  }

  return mux
}

func sortRoutes(routes map[string]http.Handler) (map[string]http.Handler, map[string]http.Handler) {
  methods, paths := make(map[string]http.Handler), make(map[string]http.Handler)

  for key, handler := range routes {
    switch key {
    case http.MethodGet, http.MethodHead, http.MethodPost,
    http.MethodPut, http.MethodPatch, http.MethodDelete,
    http.MethodConnect, http.MethodOptions, http.MethodTrace:
      methods[key] = handler
    default:
      paths[key] = handler
    }
  }

  return methods, paths
}

func defaultNotFound(res http.ResponseWriter, req *http.Request) {
  http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func defaultMethodNotAllowed(res http.ResponseWriter, req *http.Request) {
  http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
