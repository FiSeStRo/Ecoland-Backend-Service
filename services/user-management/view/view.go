package view

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// TemplateRenderer handles template rendering
type TemplateRenderer struct {
	templates map[string]*template.Template
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(templateDir string) (*TemplateRenderer, error) {
	templates := make(map[string]*template.Template)

	pages := []string{
		"users.html",
	}

	for _, page := range pages {
		tmpl, err := template.ParseFiles(
			filepath.Join(templateDir, "layout.html"),
			filepath.Join(templateDir, page),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", page, err)
		}
		templates[page] = tmpl
	}

	return &TemplateRenderer{
		templates: templates,
	}, nil
}

// Render renders a template with the provided data
func (tr *TemplateRenderer) Render(w http.ResponseWriter, name string, data any) {
	tmpl, ok := tr.templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("Template %s not found", name), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
