package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	files := []string{
		"ui/html/base.tmpl.html",
		"ui/html/pages/home.tmpl.html",
		"ui/html/partials.tmpl.html",
	}

	fs, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = fs.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return

	}
}
func (app *application) twitCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "Create a twit!")
}
func (app *application) twitView(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Fprintf(w, "Display a specific twit with an id ....%v", id)

}
