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

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	protected := dynamic.Append(app.requireAuthentication)
	//handlers
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodPost, "/twit/create", protected.ThenFunc(app.twitCreatePost))
	router.Handler(http.MethodGet, "/twit/view/:id", dynamic.ThenFunc(app.twitView))
	router.Handler(http.MethodGet, "/twit/create", protected.ThenFunc(app.twitCreate))

	//authentication and registration handlers
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogOut))

	//chaining of middleware
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//wraping it with the routers
	return standard.Then(router)
}
