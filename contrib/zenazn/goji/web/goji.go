// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

// Package web provides functions to trace the zenazn/goji/web package (https://github.com/zenazn/goji).
package web // import "gopkg.in/DataDog/dd-trace-go.v1/contrib/zenazn/goji/web"

import (
	"fmt"
	"math"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/contrib/internal/httputil"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/zenazn/goji/web"
)

// Middleware returns a goji middleware function that will trace incoming requests.
// If goji's Router middleware is also installed, the tracer will be able to determine
// the original route name (e.g. "/user/:id"), and include it as part of the traces' resource
// names.
func Middleware(opts ...Option) func(*web.C, http.Handler) http.Handler {
	var cfg config
	defaults(&cfg)
	for _, fn := range opts {
		fn(&cfg)
	}
	if !math.IsNaN(cfg.analyticsRate) {
		cfg.spanOpts = append(cfg.spanOpts, tracer.Tag(ext.EventSampleRate, cfg.analyticsRate))
	}
	return func(c *web.C, h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resource := r.Method
			p := web.GetMatch(*c).RawPattern()
			if p != nil {
				resource += " " + fmt.Sprintf("%s", p)
			}
			httputil.TraceAndServe(h, w, r, cfg.serviceName, resource, cfg.spanOpts...)
		})
	}
}
