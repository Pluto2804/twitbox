package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/twit/create", twitCreate)
	mux.HandleFunc("/twit/view", twitView)
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)

}
