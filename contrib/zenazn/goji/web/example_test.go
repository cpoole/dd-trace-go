package web_test

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	webtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/zenazn/goji/web"
)

func ExampleMiddleware() {
	// Using the Router middleware lets the tracer determine routes.
	// Otherwise the resource is "unknown".
	goji.Use(goji.DefaultMux.Router)
	goji.Use(webtrace.Middleware())
	goji.Get("/hello", func(c web.C, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Why hello there!")
	})
	goji.Serve()
}
