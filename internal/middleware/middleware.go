package middleware

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func VerboseErrorLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		if ww.Status() >= 400 {
			log.Printf("Error: %s %s - Status: %d, User-Agent: %s, Remote IP: %s\n",
				r.Method, r.RequestURI, ww.Status(), r.UserAgent(), r.RemoteAddr)
		}
	})
}
