package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"twitbox.vedantkugaonkar.net/internal/model"
)

func (app *application) home(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	twits, err := app.twits.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(req)
	data.Twits = twits
	app.renderer(w, "home.tmpl.html", http.StatusOK, data)
}

func (app *application) twitCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O shell6"
	content := "O snail\nfly through mount Himaparva,\nBut slowly,slowly!\n\n- Koba sao"
	expires := 7
	id, err := app.twits.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	//redirecting user to the relevant page for the twit
	http.Redirect(w, req, fmt.Sprintf("/twit/view?id=%d", id), http.StatusSeeOther)
}
func (app *application) twitView(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
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
