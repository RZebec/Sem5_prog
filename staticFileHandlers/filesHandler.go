package staticFileHandlers

import "net/http"

func StaticFileHandler() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/login", LoginPageHandler)
	CSSFileHandler()
	JSFileHandler()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}