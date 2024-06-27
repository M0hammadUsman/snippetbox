package main

import (
	"fmt"
	"net/http"
)

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle(fmt.Sprint(http.MethodGet, " /static/"), http.StripPrefix("/static/", fs))

	mux.HandleFunc(fmt.Sprint(http.MethodGet, " /{$}"), app.home)
	mux.HandleFunc(fmt.Sprint(http.MethodGet, " /snippets/view"), app.snippetView)
	mux.HandleFunc(fmt.Sprint(http.MethodPost, " /snippets/create"), app.snippetCreate)

	return mux

}
