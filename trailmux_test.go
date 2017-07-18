package trailmux_test

import (
	"net/http"

	"github.com/crhntr/trailmux"
)

var (
	pathHandler   http.Handler = trailmux.PathMux{}
	methodHandler http.Handler = trailmux.MethodMux{}
)
