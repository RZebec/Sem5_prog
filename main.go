package main

import (
	"net/http"
	"./staticFileHandlers"
)

func main() {
	http.HandleFunc("/", staticFileHandlers.IndexHandler)
	http.HandleFunc("/files/styles", staticFileHandlers.CssHandler)
	http.HandleFunc("/files/javascript", staticFileHandlers.JsHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}