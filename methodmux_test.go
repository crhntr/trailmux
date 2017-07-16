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

func TestIsMethod(t *testing.T) {
	for _, method := range []string{"GET", "POST", "DELETE", "PUT", "PATCH", "HEAD", "CONNECT", "OPTIONS", "TRACE"} {
		if !trailmux.IsMethod(method) {
			t.Fail()
		}
	}
	for _, method := range []string{"GETFOO", "post", "DESTROY"} {
		if trailmux.IsMethod(method) {
			t.Fail()
		}
	}
}

func TestMethodMuxInvalidMethod(t *testing.T) {
	for _, invalidMethod := range []string{"foo", "gEt", "BAR", "post"} {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("should not allow invalid method")
				}
			}()

			mux := trailmux.MethodMux{}
			mux.Handle(invalidMethod, HandlerHit{})
		}()
	}
}

func TestMethodMuxRepeatHandlerSet(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should not allow setting the same method handler twice")
		}
	}()

	target := false
	h1 := GenerateHandlerHit(&target)
	h2 := GenerateHandlerHit(&target)

	mux := trailmux.MethodMux{}
	mux.GET(h1)
	mux.GET(h2)
}

func TestMethodMuxMethodNotSet(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.GET(h)
	req, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)
	if target {
		t.Error("POST routed to incorrect Hanler")
	}
	if response.Code != http.StatusMethodNotAllowed {
		t.Error("should respond with status not found")
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
		t.Error("should use panic handler")
	}
}

// Testing Convienience Methods

func TestMethodMux_GET(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.GET(h)
	req, _ := http.NewRequest("GET", "/", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	if !target {
		t.Error("GET not handled properly")
	}
}

func TestMethodMux_POST(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.POST(h)
	req, _ := http.NewRequest("POST", "/", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	if !target {
		t.Error("POST not handled properly")
	}
}

func TestMethodMux_DELETE(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.DELETE(h)
	req, _ := http.NewRequest("DELETE", "/", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	if !target {
		t.Error("DELETE not handled properly")
	}
}

func TestMethodMux_PUT(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.PUT(h)
	req, _ := http.NewRequest("PUT", "/", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	if !target {
		t.Error("PUT not handled properly")
	}
}

func TestMethodMux_PATCH(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.PATCH(h)
	req, _ := http.NewRequest("PATCH", "/", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	if !target {
		t.Error("PATCH not handled properly")
	}
}

func TestMethodMux_HEAD(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.HEAD(h)
	req, _ := http.NewRequest("HEAD", "/", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	if !target {
		t.Error("HEAD not handled properly")
	}
}

func TestMethodMux_CONNECT(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.CONNECT(h)
	req, _ := http.NewRequest("CONNECT", "/", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	if !target {
		t.Error("CONNECT not handled properly")
	}
}

func TestMethodMux_OPTIONS(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.OPTIONS(h)
	req, _ := http.NewRequest("OPTIONS", "/", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	if !target {
		t.Error("OPTIONS not handled properly")
	}
}

func TestMethodMux_TRACE(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.MethodMux{}
	mux.TRACE(h)
	req, _ := http.NewRequest("TRACE", "/", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	if !target {
		t.Error("TRACE not handled properly")
	}
}
