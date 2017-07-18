package trailmux_test

import (
	"net/http"

	"github.com/crhntr/trailmux"
)

var (
	pathHandler   http.Handler = trailmux.PathMux{}
	methodHandler http.Handler = trailmux.MethodMux{}
)

// HandlerHit is a helper for testing that a handler selected
type HandlerHit struct {
	hit *bool
}

// GenerateHandlerHit generates a HandlerHit
func GenerateHandlerHit(target *bool) HandlerHit {
	return HandlerHit{
		hit: target,
	}
}

// ServeHTTP implements for HandlerHit
func (handler HandlerHit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	(*handler.hit) = true
}
