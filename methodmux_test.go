package trailmux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crhntr/trailmux"
)

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

func TestMethodMuxInValidMethods(t *testing.T) {
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

	for _, method := range []string{"foo", "get", "Post", "gEt"} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Error("should not panic with invalid methods")
				}
			}()
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(method, "/", nil)
			if method == "" {
				r.Method = ""
			}
			mux.ServeHTTP(w, r)

			if w.Code != http.StatusBadRequest {
				t.Errorf("method %q (%q) should return status bad request instead got status [%d] %q", method, r.Method, w.Code, http.StatusText(w.Code))
			}
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
