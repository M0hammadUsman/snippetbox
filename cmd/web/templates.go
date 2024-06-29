package main

import (
	"github.com/M0hammadUsman/snippetbox/internal/models"
	"html/template"
	"path/filepath"
	"time"
)

type templateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
	Form        any
	Flash       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{"humanDate": humanDate}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pagePaths, er := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if er != nil {
		return nil, er
	}
	for _, pp := range pagePaths {
		name := filepath.Base(pp)
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(pp)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
