package main

import (
	"fmt"
	"net/http"
)

func secureHeader(next http.Handler) http.Handler {
	x := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(x)
}

func (app *application) logRequest(next http.Handler) http.Handler {
	x := func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(x)
}

func (app *application) handlePanic(next http.Handler) http.Handler {
	x := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(x)
}
