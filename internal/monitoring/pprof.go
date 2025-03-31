package monitoring

import (
	"net/http"
	"net/http/pprof"

	"task4.2.3/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterPprofRoutes(r chi.Router) {
	r.Route("/mycustompath/pprof", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Get("/", pprof.Index)
		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/profile", pprof.Profile)
		r.Get("/symbol", pprof.Symbol)
		r.Get("/trace", pprof.Trace)

		r.Get("/allocs", func(w http.ResponseWriter, r *http.Request) {
			pprof.Handler("allocs").ServeHTTP(w, r)
		})
		r.Get("/block", func(w http.ResponseWriter, r *http.Request) {
			pprof.Handler("block").ServeHTTP(w, r)
		})
		r.Get("/goroutine", func(w http.ResponseWriter, r *http.Request) {
			pprof.Handler("goroutine").ServeHTTP(w, r)
		})
		r.Get("/heap", func(w http.ResponseWriter, r *http.Request) {
			pprof.Handler("heap").ServeHTTP(w, r)
		})
		r.Get("/mutex", func(w http.ResponseWriter, r *http.Request) {
			pprof.Handler("mutex").ServeHTTP(w, r)
		})
		r.Get("/threadcreate", func(w http.ResponseWriter, r *http.Request) {
			pprof.Handler("threadcreate").ServeHTTP(w, r)
		})
	})
}

func RegisterGoroutineDumpRoute(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Route("/mycustompath/debug/pprof", func(r chi.Router) {
			r.Get("/goroutine", func(w http.ResponseWriter, r *http.Request) {
				pprof.Handler("goroutine").ServeHTTP(w, r)
			})
		})
	})
}
