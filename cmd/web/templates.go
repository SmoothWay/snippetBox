package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/SmoothWay/snippetBox/pkg/forms"
	"github.com/SmoothWay/snippetBox/pkg/models"
)

type templateData struct {
	CurrentYear       int
	CSRFToken         string
	Flash             string
	Form              *forms.Form
	AuthenticatedUser *models.User
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*-page.html"))
	if err != nil {
		return nil, err
	}

	// loop through the pages one by one.

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*-layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*-partial.html"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts

	}
	return cache, nil
}
