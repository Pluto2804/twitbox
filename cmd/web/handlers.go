package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	ts, err := template.ParseFiles("ui/html/pages/home.tmpl.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return

	}
}
func twitCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not Allowed!", http.StatusNoContent)
		w.Header().Set("Allow", http.MethodPost)
		return
	}
	fmt.Fprintln(w, "Create a twit!")
}
func twitView(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Display a specific twit with an id ....%v", id)

}
