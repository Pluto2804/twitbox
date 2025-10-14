package main

import "net/http"

func (app *application) routeMux() http.Handler {

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/twit/create", app.twitCreate)
	mux.HandleFunc("/twit/view", app.twitView)
	return app.logRequest(secureHeaders(mux))
}
