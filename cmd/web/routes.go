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

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	//handlers
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodPost, "/twit/create", dynamic.ThenFunc(app.twitCreatePost))
	router.Handler(http.MethodGet, "/twit/view/:id", dynamic.ThenFunc(app.twitView))
	router.Handler(http.MethodGet, "/twit/create", dynamic.ThenFunc(app.twitCreate))

	//chaining of middleware
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//wraping it with the routers
	return standard.Then(router)
}
