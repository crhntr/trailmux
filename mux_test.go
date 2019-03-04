package trailmux_test

import (
  "bytes"
  "testing"
  "net/http"
  "net/http/httptest"

  "github.com/crhntr/trailmux"
)

func TestMuxHappyPaths(t *testing.T) {
  var (
    putCalled, getCalled, someOther bool
  )

  putFn := func(res http.ResponseWriter, req *http.Request) {
    putCalled = true
  }
  getFn := func(res http.ResponseWriter, req *http.Request) {
    getCalled = true
  }
  otherGetFn := func(res http.ResponseWriter, req *http.Request) {
    someOther = true
  }

  mux := trailmux.Routes{
    "/some": trailmux.Routes{
      "/path": trailmux.Routes{
        http.MethodPost: http.HandlerFunc(putFn),
      }.Mux(),
      "/other": trailmux.Routes{
        http.MethodGet: http.HandlerFunc(otherGetFn),
      }.Mux(),
    }.Mux(),
    http.MethodGet: http.HandlerFunc(getFn),
  }.Mux()

  defer func() {
    if !putCalled {
      t.Error("it should call putFn")
    }
    if !getCalled {
      t.Error("it should call getFn")
    }
    if !someOther {
      t.Error("it should call otherGetFn")
    }
  }()

  mux.ServeHTTP(
    httptest.NewRecorder(),
    httptest.NewRequest(http.MethodGet, "/", nil),
  )
  mux.ServeHTTP(
    httptest.NewRecorder(),
    httptest.NewRequest(http.MethodPost, "/some/path", bytes.NewBuffer(nil)),
  )
  mux.ServeHTTP(
    httptest.NewRecorder(),
    httptest.NewRequest(http.MethodGet, "/some/other/random/path", bytes.NewBuffer(nil)),
  )
}

func TestMuxSadPaths(t *testing.T) {
  var (
    putCalled, noMatchCalled bool
  )

  putFn := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
    putCalled = true
  })
  getFn := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {})
  noMatch := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
    res.WriteHeader(http.StatusNotFound)
    noMatchCalled = true
  })

  mux := trailmux.Routes{
    "/some": trailmux.Routes{
      "/other": trailmux.Routes{}.Mux(),
      "/path": trailmux.Routes{
        http.MethodPost: putFn,
      }.Mux(),
    }.Mux(),
    http.MethodGet: getFn,
    "/no-match": trailmux.Routes{}.Mux().NoMatchHandler(noMatch),
  }.Mux()

  defer func() {
    if putCalled {
      t.Error("it should not call putFn")
    }
  }()

  {
    res := httptest.NewRecorder()
    mux.ServeHTTP(
      res,
      httptest.NewRequest(http.MethodDelete, "/", nil),
    )
    if res.Code != http.StatusNotFound {
      t.Errorf("expected %d got %d", http.StatusNotFound, res.Code)
    }
  }

  {
    res := httptest.NewRecorder()
    mux.ServeHTTP(
      res,
      httptest.NewRequest(http.MethodDelete, "/some/path", nil),
    )
    if res.Code != http.StatusMethodNotAllowed {
      t.Errorf("expected %d got %d", http.StatusMethodNotAllowed, res.Code)
    }
  }

  {
    res := httptest.NewRecorder()
    mux.ServeHTTP(
      res,
      httptest.NewRequest(http.MethodGet, "/some/other", nil),
    )
    if res.Code != http.StatusNotFound {
      t.Errorf("expected %d got %d", http.StatusNotFound, res.Code)
    }
  }

  {
    res := httptest.NewRecorder()
    mux.ServeHTTP(
      res,
      httptest.NewRequest(http.MethodGet, "/no-match", nil),
    )
    if res.Code != http.StatusNotFound {
      t.Errorf("expected %d got %d", http.StatusNotFound, res.Code)
    }
    if !noMatchCalled {
      t.Error("it should not call noMatch")
    }
  }
}
