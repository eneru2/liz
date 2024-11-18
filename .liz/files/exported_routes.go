package files

import (
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func Routes(r *chi.Mux) {
	r.Get("/hello",templ.Handler(Route_hello()).ServeHTTP)
r.Get("/hello/hey",templ.Handler(Route_hello_hey()).ServeHTTP)
r.Get("/", templ.Handler(Route_root()).ServeHTTP)

}
