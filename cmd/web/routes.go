package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)


func (app *application) routes() http.Handler {
    router := httprouter.New()

    router.NotFound = app.routerWrap()
    router.MethodNotAllowed = app.routerWrap()

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

    // base
    standard := alice.New(
        app.recoverPanic,
        app.logRequest,
        secureHeaders,
    )

    // with session + csrf + auth check
    dynamic := standard.Append(
        app.sessionManager.LoadAndSave,
        noSurf,
        app.authenticate,
    )

    protected := dynamic.Append(app.requireAuthentication)

    // public
    router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
    router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))
    router.Handler(http.MethodGet, "/twit/view/:id", dynamic.ThenFunc(app.twitView))
    router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
    router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
    router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
    router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

    // protected
    router.Handler(http.MethodGet, "/twit/create", protected.ThenFunc(app.twitCreate))
    router.Handler(http.MethodPost, "/twit/create", protected.ThenFunc(app.twitCreatePost))
    router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogOut))

    return dynamic.Then(router) // ✅ not standard.Then
}
