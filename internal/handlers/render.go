package handlers

import (
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/ironedr/gochat/ui"
)

var errorTemplateNotFound = errors.New("template not found")

type templateData struct {
	Form            any
	StringMap       map[string]string
	CurrentYear     int
	IsAuthenticated bool
}

// getDate returns the formatted date
func getDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// isModuloZero returns true if modulo operation between the two input parameters returns 0
func isModuloZero(val, mod int) bool {
	return val%mod == 0
}

var templateFuncs = template.FuncMap{
	"getDate": getDate,
	"isMod":   isModuloZero,
}

func (h *Handler) NewTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// NewTemplateCache iterates through the html templates in the ui package and caches them, returning a map
func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	var pageDependencies []string
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	layout, err := fs.Glob(ui.Files, "html/base.tmpl")
	if err != nil {
		return nil, err
	}
	pageDependencies = append(pageDependencies, layout[0])

	partials, err := fs.Glob(ui.Files, "html/partials/*.tmpl")
	if err != nil {
		return nil, err
	}

	pageDependencies = append(pageDependencies, partials...)

	for _, page := range pages {
		pageName := filepath.Base(page)
		ts, err := template.New(pageName).Funcs(templateFuncs).ParseFS(ui.Files, page)
		if err != nil {
			return nil, err
		}
		for _, pageDependency := range pageDependencies {
			ts, err = ts.ParseFS(ui.Files, pageDependency)
			if err != nil {
				return nil, err
			}
		}
		// so that we can simply use the filename without the extension as our key
		cache[strings.TrimSuffix(pageName, filepath.Ext(page))] = ts
	}
	return cache, nil
}

// Render a given template if it exists in the template cache
func (h *Handler) Render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := h.app.TemplateCache[page]
	if !ok {
		h.ServerError(w, errorTemplateNotFound)
		return
	}
	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.ServerError(w, err)
	}
}
