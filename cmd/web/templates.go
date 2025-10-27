package main

import (
	"path/filepath"
	"text/template"
	"time"

	"twitbox.vedantkugaonkar.net/internal/model"
)

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type templateData struct {
	CurrentYear     int
	Twit            *model.Twit
	Twits           []*model.Twit
	Form            any
	Flash           string
	IsAuthenticated bool
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	//use glob function to get the slice of all filepaths that
	//match the base provided
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		//extract the base name of the filepath for cache map
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles("ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("ui/html/partials*.tmpl.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil

}
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
