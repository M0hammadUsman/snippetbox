package main

import (
	"bytes"
	"errors"
	"github.com/go-playground/form/v4"
	"log/slog"
	"net/http"
	"time"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Error("the template %s does not exist", "page", page)
		return
	}
	b := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(b, "base", data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	w.WriteHeader(status)
	_, err := b.WriteTo(w)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Error(err.Error())
	}
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	if err := app.formDecoder.Decode(dst, r.PostForm); err != nil {
		var invalidDecoderErr *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderErr) {
			panic(err)
		}
		return err
	}
	return nil
}
