package main

import (
	"path/filepath"
	"text/template"

	"twitbox.vedantkugaonkar.net/internal/model"
)

type templateData struct {
	Twit  *model.Twit
	Twits []*model.Twit
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
		//slice of paths for the base template
		files := []string{
			"ui/html/base.tmpl.html",
			"ui/html/partials.tmpl.html",
			page,
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil

}
