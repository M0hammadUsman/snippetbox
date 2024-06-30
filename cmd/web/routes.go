package main

import (
	"fmt"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle(fmt.Sprint(http.MethodGet, " /static/"), http.StripPrefix("/static/", fs))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle(fmt.Sprint(http.MethodGet, " /{$}"), dynamic.ThenFunc(app.home))
	mux.Handle(fmt.Sprint(http.MethodGet, " /snippets/view/{id}"), dynamic.ThenFunc(app.snippetView))
	mux.Handle(fmt.Sprint(http.MethodGet, " /snippets/create"), dynamic.ThenFunc(app.snippetCreate))
	mux.Handle(fmt.Sprint(http.MethodPost, " /snippets/create"), dynamic.ThenFunc(app.snippetCreatePost))

	mux.Handle(fmt.Sprint(http.MethodGet, " /users/signup"), dynamic.ThenFunc(app.userSignup))
	mux.Handle(fmt.Sprint(http.MethodPost, " /users/signup"), dynamic.ThenFunc(app.userSignupPost))
	mux.Handle(fmt.Sprint(http.MethodGet, " /users/login"), dynamic.ThenFunc(app.userLogin))
	mux.Handle(fmt.Sprint(http.MethodPost, " /users/login"), dynamic.ThenFunc(app.userLoginPost))
	mux.Handle(fmt.Sprint(http.MethodPost, " /users/logout"), dynamic.ThenFunc(app.userLogoutPost))

	standard := alice.New(recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
}
