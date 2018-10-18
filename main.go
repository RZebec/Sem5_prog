package main

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	var message =
		"<html>" +
			"<body>" +
			"<h1>Welcome to the Index Page</h1>" +
			"</body>" +
			"</html>"

	w.Write([]byte(message))
}

func main() {
	http.HandleFunc("/", index)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
