package main

import (
	"fmt"
	"net/http"
)

/*
Middleware basically is a function that is called at the beginning of a handler.
The difference is it outputs a request and middleware can be chained.
It would be incredibly cumbersome and inefficient to implement the same functions for every single endpoint, especially when you have a lot of them.
It’s just a function that processes request so your handler can focus on one thing: handling the request after it’s already been processed.
You don’t need to check if the request is authenticated because your authentication middleware has already done that.
Middleware gives you the ability to begin your handler already operating with certain assumptions in place.
If you have 10 different middleware functions, it’s certainly not DRY-compliant to write those same 10 function calls at the beginning of every handler. You could, but why would you want to?
*/
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self' font.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, req)

	})
}
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", req.RemoteAddr, req.Proto, req.Method, req.URL.RequestURI())
		next.ServeHTTP(w, req)
	})
}
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, req)
	})
}
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !app.isAuthenticated(req) {
			http.Redirect(w, req, "/user/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, req)
	})
}
