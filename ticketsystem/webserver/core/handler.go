package core

import (
	"fmt"
	"net/http"
	"time"
)

type HttpHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	Next HttpHandler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Do stuff before execution of child handler
	fmt.Println("HttpHandler Before" + time.Now().String())
	defer fmt.Println("HttpHandler After")
	h.Next.ServeHTTP(w, r)
	// Do stuff after execution of child handler

}
