package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"twitbox.vedantkugaonkar.net/internal/model"
	"twitbox.vedantkugaonkar.net/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
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
	//when httprouter parses a req, any named parameters are stored in the
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
	var scF snippetCreateForm
	err := app.decodePostForm(req, &scF)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	scF.CheckField(validator.NotBlank(scF.Title), "title", "This field cannot be blank!")
	scF.CheckField(validator.MaxChars(scF.Title, 100), "title", "This field cannot be more than 100 characters")
	scF.CheckField(validator.PermittedInt(scF.Expires, 1, 7, 365), "expires", "This field must include 1,7 or 365")
	scF.CheckField(validator.NotBlank(scF.Content), "content", "This field cannot be empty!")
	if !scF.Valid() {
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
