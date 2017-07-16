package trailmux

// Some of the following code was copied from a repo with the following notice
// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// at https://github.com/julienschmidt/httprouter/blob/master/LICENSE

import (
	"net/http"
)

// PathMux is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type PathMux struct {
	root *node

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound http.Handler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})

	RedirectTrailingSlash bool
}

// Make sure the PathMux conforms with the http.Handler interface
var _ http.Handler = &PathMux{}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *PathMux) Handle(path string, handler http.Handler) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.root == nil {
		r.root = &node{}
	}

	r.root.addRoute(path, handler)
}

func (r *PathMux) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *PathMux) Lookup(path string) (http.Handler, Params, bool) {
	return r.root.getValue(path)
}

// func (r *PathMux) allowed(path, reqMethod string) (allow string) {
// 	if path == "*" { // server-wide
// 		for method := range r.trees {
// 			if method == "OPTIONS" {
// 				continue
// 			}
//
// 			// add request method to list of allowed methods
// 			if len(allow) == 0 {
// 				allow = method
// 			} else {
// 				allow += ", " + method
// 			}
// 		}
// 	} else { // specific path
// 		for method := range r.trees {
// 			// Skip the requested method - we already tried this one
// 			if method == reqMethod || method == "OPTIONS" {
// 				continue
// 			}
//
// 			handle, _, _ := r.trees[method].getValue(path)
// 			if handle != nil {
// 				// add request method to list of allowed methods
// 				if len(allow) == 0 {
// 					allow = method
// 				} else {
// 					allow += ", " + method
// 				}
// 			}
// 		}
// 	}
// 	if len(allow) > 0 {
// 		allow += ", OPTIONS"
// 	}
// 	return
// }

// ServeHTTP makes the router implement the http.Handler interface.
func (r *PathMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}

	path := req.URL.Path

	if r.root != nil {
		if handle, ps, tsr := r.root.getValue(path); handle != nil {
			handle.ServeHTTP(w, setParamsInContext(req, ps))
			return
		} else if req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}
