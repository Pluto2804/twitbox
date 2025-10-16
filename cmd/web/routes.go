package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routeMux() http.Handler {
	router := httprouter.New()

	router.NotFound = app.routerWrap()
	router.MethodNotAllowed = app.routerWrapA()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	//handlers
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodPost, "/twit/create", app.twitCreatePost)
	router.HandlerFunc(http.MethodGet, "/twit/view/:id", app.twitView)
	router.HandlerFunc(http.MethodGet, "/twit/create", app.twitCreate)

	//chaining of middleware
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//wraping it with the routers
	return standard.Then(router)
}
