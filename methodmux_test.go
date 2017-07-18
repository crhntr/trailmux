package trailmux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crhntr/trailmux"
)

// HandlerHit is a helper for testing that a handler selected
type HandlerHit struct {
	hit *bool
}

func GenerateHandlerHit(target *bool) HandlerHit {
	return HandlerHit{
		hit: target,
	}
}

func (handler HandlerHit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	(*handler.hit) = true
}

var methodStrings = [...]string{"GET", "POST", "DELETE", "PUT", "PATCH", "HEAD", "CONNECT", "OPTIONS", "TRACE"}

func TestMethodMuxValidMethods(t *testing.T) {
	target := [len(methodStrings)]bool{}

	mux := trailmux.MethodMux{
		GET:     GenerateHandlerHit(&target[0]),
		POST:    GenerateHandlerHit(&target[1]),
		DELETE:  GenerateHandlerHit(&target[2]),
		PUT:     GenerateHandlerHit(&target[3]),
		PATCH:   GenerateHandlerHit(&target[4]),
		HEAD:    GenerateHandlerHit(&target[5]),
		CONNECT: GenerateHandlerHit(&target[6]),
		OPTIONS: GenerateHandlerHit(&target[7]),
		TRACE:   GenerateHandlerHit(&target[8]),
	}

	for _, method := range methodStrings {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Error("should allow valid methods")
				}
			}()
			w := httptest.NewRecorder()
			r, err := http.NewRequest(method, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			mux.ServeHTTP(w, r)
		}()
	}
}

func TestMethodMuxMethodNotSet(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.GET = h
	req, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)
	if target {
		t.Error("POST routed to incorrect Hanler")
	}
	if response.Code != http.StatusMethodNotAllowed {
		t.Error("should respond with status method not allowed")
	}
}

func TestMethodMuxMethodNotFoundUse(t *testing.T) {
	target := false
	notFoundHandler := GenerateHandlerHit(&target)

	mux := trailmux.MethodMux{}
	mux.MethodNotAllowed = notFoundHandler
	req, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)

	if !target {
		t.Error("should use method not found handler")
	}
}
