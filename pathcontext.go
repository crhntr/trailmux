package trailmux

// Some of the following code was copied from a repo with the following notice
// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// at https://github.com/julienschmidt/httprouter/blob/master/LICENSE

import (
	"context"
	"net/http"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) Get(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

type paramsIDKeyType uint8

const paramsIDKey paramsIDKeyType = 42

func setParamsInContext(r *http.Request, ps Params) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), paramsIDKey, ps))
}

func PathParameters(r http.Request) (params Params) {
	params, _ = r.Context().Value(paramsIDKey).(Params)
	return params
}
