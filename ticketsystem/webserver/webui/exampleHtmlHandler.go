package webui

import (
	"fmt"
	"net/http"
	"strings"
)


type ExampleHtmlHandler struct {
	Prefix string
	}


func (e ExampleHtmlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Println(e.Prefix + " ExampleHtmlHandler: Received Request, Method:" + r.Method)
	switch method := strings.ToLower(r.Method); method {
	case "post":
		handlePost(w,r)
	case "put":
		HandlePut(w,r)
	case "get":
		HandleGet(w,r)
	case "delete":
		HandleDelete(w,r)
	case "patch":
		HandlePatch(w,r)
	}
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Good"))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func  HandlePatch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func  HandleDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func  HandlePut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}