/*
Helpful links to read up on go middlewares

* https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81#.oqdtvyupk
* https://justinas.org/writing-http-middleware-in-go/
* https://www.youtube.com/watch?v=xyDkyFjzFVc

*/

package middlewares

import (
	"log"
	"net/http"
	"strings"
)

// Middleware is a type for decorating requests
type Middleware func(http.Handler) http.Handler

// Apply a middlewares to a handler
func Apply(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, adapter := range middlewares {
		h = adapter(h)
	}
	return h
}

// Logging is a middleware for adding a request log
func Logging() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := []string{
				r.URL.Path,
				r.Referer(),
				r.Header.Get("User-Agent"),
				r.RemoteAddr,
			}
			log.Println(strings.Join(parts, " - "))
			h.ServeHTTP(w, r)
		})
	}
}
