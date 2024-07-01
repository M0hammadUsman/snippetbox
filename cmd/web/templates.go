package main

import (
	"github.com/M0hammadUsman/snippetbox/internal/models"
	"github.com/M0hammadUsman/snippetbox/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

type templateData struct {
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{"humanDate": humanDate}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pagePaths, er := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if er != nil {
		return nil, er
	}
	for _, pp := range pagePaths {
		name := filepath.Base(pp)
		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			pp,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
