package trailmux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crhntr/trailmux"
)

func TestPathMuxStaticPath(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.PathMux{}
	mux.Handle("/foo", h)
	req, _ := http.NewRequest("GET", "/foo", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)
	if !target {
		t.Error("path should be found")
	}
	if response.Code != http.StatusOK {
		t.Error("should respond with status ok")
	}
}

func TestPathMuxGetParams(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := trailmux.PathParameters(r)
		if params.Get("msg") != "helloworld" {
			t.Error("should give valid params")
		}
		if params.Get("othermsg") != "" {
			t.Error("should give empty param for non existant param")
		}
	})
	mux := trailmux.PathMux{}
	mux.Handle("/foo/:msg", h)
	req, _ := http.NewRequest("GET", "/foo/helloworld", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Error("should respond with status ok")
	}
}

func TestPathMuxPathNotFound(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.PathMux{}
	mux.Handle("/foo", h)
	req, _ := http.NewRequest("GET", "/bar", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)
	if target {
		t.Error("handler should not be reached")
	}
	if response.Code != http.StatusNotFound {
		t.Error("should respond with status not found")
	}
}

func TestPathMuxPathParam(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.PathMux{}
	mux.Handle("/path/:var", h)
	req, _ := http.NewRequest("GET", "/path/foo", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)
	if !target {
		t.Error("handler should be reached")
	}
	if response.Code != http.StatusOK {
		t.Error("should respond with status ok")
	}
}

func TestPathMuxDynamicAndStaticPath(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.PathMux{}
	mux.Handle("/path/:var1/route/:var2", h)
	req, _ := http.NewRequest("GET", "/path/foo/route/bar", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)
	if !target {
		t.Error("handler should be reached")
	}
	if response.Code != http.StatusOK {
		t.Error("should respond with status ok")
	}
}

func TestPathMuxEmptyPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should not allow setting an empty path")
		}
	}()

	target := false
	h := GenerateHandlerHit(&target)

	mux := trailmux.PathMux{}
	mux.Handle("", h)
}

func TestPathMuxPathWithoutStartingSlash(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should not allow setting an empty path")
		}
	}()

	target := false
	h := GenerateHandlerHit(&target)

	mux := trailmux.PathMux{}
	mux.Handle("hello", h)
}

func TestPathMuxPanicHandlerUse(t *testing.T) {
	target := false
	panicHandler := func(w http.ResponseWriter, r *http.Request, resp interface{}) {
		target = true
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("why?")
	})

	mux := trailmux.PathMux{}
	mux.PanicHandler = panicHandler
	mux.Handle("/panic", h)

	req, _ := http.NewRequest("GET", "/panic", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)

	if !target {
		t.Error("should use panic handler")
	}
}

func TestPathMuxNotFoundHandlerUse(t *testing.T) {
	target := false
	notFoundTarget := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target = true
	})

	mux := trailmux.PathMux{}
	mux.NotFound = notFoundTarget

	req, _ := http.NewRequest("GET", "/somepath", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)

	if !target {
		t.Error("should use not found handler")
	}
}

func TestPathMuxAddingSimularStaticPaths(t *testing.T) {
	paths := [...]string{
		"/foo",
		"/bar",
		"/foo/bar",
		"/f/b",
		"/foo/bar/baz",
		"/foo/bar/baz/",
	}

	targets := [len(paths)]bool{}
	for i, _ := range targets {
		targets[i] = false
	}

	mux := trailmux.PathMux{}

	for _, path := range paths {
		mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for pathIndex, path := range paths {
				if path == r.URL.Path {
					targets[pathIndex] = true
				}
			}
		}))
	}

	for i, path := range paths {
		req, _ := http.NewRequest("GET", path, nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, req)

		if !targets[i] || response.Code != http.StatusOK {
			t.Errorf("path not reached %s", path)
		}
	}
}

func TestPathMuxAddingSimularDynamicPaths(t *testing.T) {
	paths := [...]string{
		"/foo",
		"/foo/",
		"/foo/:p1",
		"/foo/:p1/",
		"/foo/:p1/bar",
		"/foo/:p1/baz",
		"/foo/:p1/baz/:p2",
		"/foo/:p1/baz/:p2/:p3",
		"/foo/:p1/baz/:p2/:p3/",
		"/foo/:p1/baz/:p2/:p3/:p4",
		"/",
	}

	targets := [len(paths)]bool{}
	for i, _ := range targets {
		targets[i] = false
	}

	mux := trailmux.PathMux{}

	for _, path := range paths {
		mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for pathIndex, path := range paths {
				if path == r.URL.Path {
					targets[pathIndex] = true
				}
			}
		}))
	}

	for i, path := range paths {
		req, _ := http.NewRequest("GET", path, nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, req)

		if !targets[i] || response.Code != http.StatusOK {
			t.Errorf("path not reached %s", path)
		}
	}
}

func TestPathMuxConflictingPaths(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should not allow handling conflicting paths")
		}
	}()
	paths := [...]string{
		"/foo",
		"/:p1",
	}

	targets := [len(paths)]bool{}
	for i, _ := range targets {
		targets[i] = false
	}

	mux := trailmux.PathMux{}

	for _, path := range paths {
		mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for pathIndex, path := range paths {
				if path == r.URL.Path {
					targets[pathIndex] = true
				}
			}
		}))
	}
}

func TestPathMuxConflictingWildCardPaths(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should not allow handling conflicting paths")
		}
	}()
	paths := [...]string{
		"/:p2",
		"/:p1",
	}

	targets := [len(paths)]bool{}
	for i, _ := range targets {
		targets[i] = false
	}

	mux := trailmux.PathMux{}

	for _, path := range paths {
		mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for pathIndex, path := range paths {
				if path == r.URL.Path {
					targets[pathIndex] = true
				}
			}
		}))
	}
}

func TestPathMuxMaxParamsExceeded(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("should not allow handling conflicting paths")
		}
	}()
	pathHandlerRoute := ""
	pathRequestRoute := ""
	for i := 0; i < trailmux.MaxParams+2; i++ {
		pathHandlerRoute += fmt.Sprintf("/:p%d", i)
		pathRequestRoute += "/p"
	}

	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.PathMux{}
	mux.Handle(pathHandlerRoute, h)
	req, _ := http.NewRequest("GET", pathRequestRoute, nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)
	if target {
		t.Error("handler should not be reached")
	}
	if response.Code != http.StatusNotFound {
		t.Error("should respond with status ok")
	}
}

func TestPathMuxShouldNotAllowExactStaticPaths(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should not allow handling conflicting paths")
		}
	}()
	paths := [...]string{
		"/foo",
		"/foo",
	}

	targets := [len(paths)]bool{}
	for i, _ := range targets {
		targets[i] = false
	}

	mux := trailmux.PathMux{}

	for _, path := range paths {
		mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for pathIndex, path := range paths {
				if path == r.URL.Path {
					targets[pathIndex] = true
				}
			}
		}))
	}
}

func TestPathMuxShouldNotAllowMultipleWildcardsInPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should not allow multiple wildcards paths")
		}
	}()
	paths := [...]string{
		"/foo/*var/*",
	}

	targets := [len(paths)]bool{}
	for i, _ := range targets {
		targets[i] = false
	}

	mux := trailmux.PathMux{}

	for _, path := range paths {
		mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for pathIndex, path := range paths {
				if path == r.URL.Path {
					targets[pathIndex] = true
				}
			}
		}))
	}
}

func TestPathMuxShouldHaveValidname(t *testing.T) {
	paths := [...]string{
		"/foo/*a*/",
		"/foo/*a:/",
		"/foo/*:/",
		"/foo/::/",
		"/foo/:a:/",
		"/:/",
		"/*/",
		"/foo*var",
	}

	for _, path := range paths {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("should not allow wildcard char in wildcardname")
				}
			}()

			targets := [len(paths)]bool{}
			for i, _ := range targets {
				targets[i] = false
			}

			mux := trailmux.PathMux{}

			mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for pathIndex, path := range paths {
					if path == r.URL.Path {
						targets[pathIndex] = true
					}
				}
			}))
		}()
	}
}

func TestPathMuxShouldNotAllowInvalidCatchAllPath(t *testing.T) {
	paths := [...]string{
		"*",
		"/*/foo",
	}

	for _, path := range paths {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("should not allow invalid path %q", path)
				}
			}()

			targets := [len(paths)]bool{}
			for i, _ := range targets {
				targets[i] = false
			}

			mux := trailmux.PathMux{}

			mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for pathIndex, path := range paths {
					if path == r.URL.Path {
						targets[pathIndex] = true
					}
				}
			}))
		}()
	}
}

func TestPathMuxShouldNotAllowCoflictingWildcard(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should not allow conflicting wildcard")
		}
	}()
	paths := [...]string{
		"/foo*var",
	}

	targets := [len(paths)]bool{}
	for i, _ := range targets {
		targets[i] = false
	}

	mux := trailmux.PathMux{}

	for _, path := range paths {
		mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for pathIndex, path := range paths {
				if path == r.URL.Path {
					targets[pathIndex] = true
				}
			}
		}))
	}
}

func TestPathMuxCatchAllPath(t *testing.T) {
	target := false
	h := GenerateHandlerHit(&target)
	mux := trailmux.PathMux{}
	mux.Handle("/foo/*var", h)
	req, _ := http.NewRequest("GET", "/foo/somefile.txt", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, req)
	if !target {
		t.Error("path should be found")
	}
	if response.Code != http.StatusOK {
		t.Error("should respond with status ok")
	}
}

func TestPathMuxLookup(t *testing.T) {
	paths := [...]string{
		"/foo",
		"/bar",
		"/baz",
		"/users/:id",
		"/endpoint/:number/:variation",
	}

	targets := [len(paths)]bool{}
	mux := trailmux.PathMux{}

	for _, path := range paths {
		mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for pathIndex, path := range paths {
				if path == r.URL.Path {
					targets[pathIndex] = true
				}
			}
		}))
	}

	for _, path := range paths {
		handler, _, _ := mux.Lookup(path)
		if handler == nil {
			t.Error("should find handler")
		}
	}
}

func TestPathMuxDisallowCatchAllConflict(t *testing.T) {
	// tests 				panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")

	defer func() {
		if r := recover(); r == nil {
			t.Error("should not allow existing handle for the path segment root in path")
		}
	}()

	target1, target2 := false, false
	h1 := GenerateHandlerHit(&target1)
	h2 := GenerateHandlerHit(&target2)

	mux := trailmux.PathMux{}
	mux.Handle("/hello/foo/", h1)
	mux.Handle("/hello/:foo/", h2)
}

func TestPathMuxRedirectRecomendationsStaticPaths(t *testing.T) {
	{
		target1, target2 := false, false
		h1 := GenerateHandlerHit(&target1)
		h2 := GenerateHandlerHit(&target2)

		mux := trailmux.PathMux{}
		mux.RedirectTrailingSlash = true
		mux.Handle("/hello/foo/", h1)
		mux.Handle("/hello/bar", h2)

		req, _ := http.NewRequest("GET", "/hello/foo", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, req)
		if target1 {
			t.Error("path should not be found")
		}
		if response.Code != http.StatusMovedPermanently {
			t.Error("should respond with status moved permanently")
		}
	}

	{
		target1, target2 := false, false
		h1 := GenerateHandlerHit(&target1)
		h2 := GenerateHandlerHit(&target2)

		mux := trailmux.PathMux{}
		mux.RedirectTrailingSlash = true
		mux.Handle("/hello/foo", h1)
		mux.Handle("/hello/bar", h2)

		req, _ := http.NewRequest("GET", "/hello/foo/", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, req)
		if target1 {
			t.Error("path should not be found")
		}
		if response.Code != http.StatusMovedPermanently {
			t.Error("should respond with status moved permanently")
		}
	}

	{
		// tests 				panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
		target1, target2 := false, false
		h1 := GenerateHandlerHit(&target1)
		h2 := GenerateHandlerHit(&target2)

		mux := trailmux.PathMux{}
		mux.RedirectTrailingSlash = true
		mux.Handle("/hello/foo/", h1)
		mux.Handle("/hello/bar", h2)

		req, _ := http.NewRequest("POST", "/hello/foo", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, req)
		if target1 {
			t.Error("path should not be found")
		}
		if response.Code != http.StatusTemporaryRedirect {
			t.Error("should respond with status temporary redirect")
		}
	}

	{
		// tests 				panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
		target1, target2 := false, false
		h1 := GenerateHandlerHit(&target1)
		h2 := GenerateHandlerHit(&target2)

		mux := trailmux.PathMux{}
		mux.RedirectTrailingSlash = true
		mux.Handle("/hello/foo", h1)
		mux.Handle("/hello/bar", h2)

		req, _ := http.NewRequest("POST", "/hello/foo/", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, req)
		if target1 {
			t.Error("path should not be found")
		}
		if response.Code != http.StatusTemporaryRedirect {
			t.Error("should respond with status temporary redirect")
		}
	}
}

func TestPathMuxRedirectRecomendationsDynamicPaths(t *testing.T) {
	{
		target1 := false
		h1 := GenerateHandlerHit(&target1)

		mux := trailmux.PathMux{}
		mux.RedirectTrailingSlash = true
		mux.Handle("/hello/:foo/", h1)

		req, _ := http.NewRequest("GET", "/hello/foo", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, req)
		if target1 {
			t.Error("path should not be found")
		}
		if response.Code != http.StatusMovedPermanently {
			t.Error("should respond with status moved permanently")
		}
	}

	{
		// tests 				panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
		target1 := false
		h1 := GenerateHandlerHit(&target1)

		mux := trailmux.PathMux{}
		mux.RedirectTrailingSlash = true
		mux.Handle("/hello/foo/", h1)

		req, _ := http.NewRequest("POST", "/hello/foo", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, req)
		if target1 {
			t.Error("path should not be found")
		}
		if response.Code != http.StatusTemporaryRedirect {
			t.Error("should respond with status temporary redirect")
		}
	}
}
