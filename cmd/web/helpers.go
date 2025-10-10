package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// sends 500 internal server error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// client error for sending a specific response to the user(i.e bad req)
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)

}

// not found response wrapper around clientError to send a 404 not found response
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
func (app *application) renderer(w http.ResponseWriter, page string, status int, data *templateData) {
	ts, ok := app.tempCache[page]
	if !ok {
		err := fmt.Errorf("template for %s not available", page)
		app.serverError(w, err)
		return
	}
	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}
