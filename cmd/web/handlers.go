package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"twitbox.vedantkugaonkar.net/internal/model"
)

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, req *http.Request) {
	twits, err := app.twits.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(req)
	data.Twits = twits
	app.renderer(w, "home.tmpl.html", http.StatusOK, data)
}
func (app *application) twitView(w http.ResponseWriter, req *http.Request) {
	//when httprouter parses a req, any named parameter are stored in the
	//req context so here ParseFromContext is used to retrieve/get the id
	params := httprouter.ParamsFromContext(req.Context())
	//ByName can then be used to get the actual value
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	twit, err := app.twits.Get(id)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(req)
	data.Twit = twit
	app.renderer(w, "view.tmpl.html", http.StatusOK, data)

}
func (app *application) twitCreate(w http.ResponseWriter, req *http.Request) {
	data := &templateData{}
	data.Form = &snippetCreateForm{
		Expires: 365,
	}
	app.renderer(w, "create.tmpl.html", http.StatusCreated, data)

}
func (app *application) twitCreatePost(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := req.PostForm.Get("title")
	content := req.PostForm.Get("Content")

	expires, err := strconv.Atoi(req.PostForm.Get("Expires"))
	//instance of snippetCreateForm struct containing the values from the form
	//along with empty map initialised for any validation errors

	scF := &snippetCreateForm{
		Title:       title,
		Content:     content,
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(title) == "" {
		scF.FieldErrors["title"] = "This field cannot be empty!"
	} else if utf8.RuneCountInString(title) > 100 {
		scF.FieldErrors["title"] = "This field cannot be more than 100 characters"
	}

	if strings.TrimSpace(content) == "" {
		scF.FieldErrors["content"] = "This field cannot be displayed!"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		scF.FieldErrors["expires"] = "This field must equal 1,7 or 365!"
	}
	if len(scF.FieldErrors) > 0 {
		data := app.newTemplateData(req)
		data.Form = scF
		app.renderer(w, "create.tmpl.html", http.StatusUnprocessableEntity, data)
		return
	}

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	id, err := app.twits.Insert(scF.Title, scF.Content, scF.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	//redirecting user to the relevant page for the twit
	http.Redirect(w, req, fmt.Sprintf("/twit/view/%d", id), http.StatusSeeOther)
}
