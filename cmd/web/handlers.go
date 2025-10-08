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
	// files := []string{
	// 	"ui/html/base.tmpl.html",
	// 	"ui/html/pages/home.tmpl.html",
	// 	"ui/html/partials.tmpl.html",
	// }

	// fs, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	// err = fs.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	twits, err := app.twits.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, twit := range twits {
		fmt.Fprintf(w, "%+v\n", twit)
	}

}
func (app *application) twitCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb mount Fuji,\nBut slowly,slowly!\n\n- Kobayashi Issa"
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
	fmt.Fprintf(w, "%+v", twit)

}
