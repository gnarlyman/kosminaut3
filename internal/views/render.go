package views

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
)

type Renderer struct {
	pages    *template.Template
	partials *template.Template
}

func New(templatesFS fs.FS) (*Renderer, error) {
	pages, err := template.ParseFS(templatesFS, "templates/layout.html", "templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("parse pages: %w", err)
	}
	partials, err := template.ParseFS(templatesFS, "templates/partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse partials: %w", err)
	}
	return &Renderer{pages: pages, partials: partials}, nil
}

func (r *Renderer) Page(w io.Writer, data any) error {
	return r.pages.ExecuteTemplate(w, "layout", data)
}

func (r *Renderer) Partial(w io.Writer, name string, data any) error {
	return r.partials.ExecuteTemplate(w, name, data)
}
