package main

import (
	"errors"
	"fmt"
	"github.com/M0hammadUsman/snippetbox/internal/models"
	"log/slog"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, _ *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	/*tf := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}
	ts, err := template.ParseFiles(tf...)
	if err != nil {
		slog.Error(http.StatusText(http.StatusInternalServerError), "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err = ts.ExecuteTemplate(w, "base", nil); err != nil {
		slog.Error(http.StatusText(http.StatusInternalServerError), "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}*/
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	res, err := app.snippets.Get(id)
	if errors.Is(err, models.ErrorNoRows) {
		http.NotFound(w, r)
		return
	}
	_, err = fmt.Fprintf(w, "%+v", res)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	err := app.snippets.Insert(title, content, expires)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippets/view?id=1"), http.StatusSeeOther)
}
