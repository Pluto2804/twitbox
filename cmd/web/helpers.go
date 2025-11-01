package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
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
func (app *application) renderer(w http.ResponseWriter,req *http.Request, page string, status int, data *templateData) {
	ts, ok := app.tempCache[page]
	if !ok {
		err := fmt.Errorf("template for %s not available", page)
		app.serverError(w, err)
		return
	}
	if data == nil {
		data = app.newTemplateData(req)
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)

}
func (app *application) newTemplateData(req *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(req.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(req),
		CSRFToken:       nosurf.Token(req),
	}

}
func (app *application) routerWrap() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		app.notFound(w)
	})
}
func (app *application) routerWrapA() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	})
}
func (app *application) decodePostForm(req *http.Request, dst any) error {
	err := req.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, req.PostForm)

	if err != nil {
		var invalidDecoderError *form.InvalidEncodeError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}
func (app *application) isAuthenticated(req *http.Request) bool {
	isAuthenticated, ok := req.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false

	}
	return isAuthenticated
}
