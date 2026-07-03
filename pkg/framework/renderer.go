package framework

import (
	"bytes"
	"net/http"
	"text/template"
)

type RendererConfig struct {
	LayoutGlob   string
	TemplateGlob string
}

type Renderer struct {
	templates *template.Template
}

type viewData struct {
	Title       string
	Description string
	CartCount   int
	Data        any
}

func NewRenderer(config RendererConfig) (*Renderer, error) {
	r := &Renderer{}

	t, err := template.New("app").ParseGlob(config.LayoutGlob)
	if err != nil {
		return nil, err
	}
	t, err = t.ParseGlob(config.TemplateGlob)
	if err != nil {
		return nil, err
	}
	r.templates = t
	return r, nil
}

func (r *Renderer) Render(w http.ResponseWriter, page Page) {
	if page.Status == 0 {
		page.Status = http.StatusOK
	}

	data := viewData{
		Title:       page.Title,
		Description: page.Description,
		Data:        page.Data,
	}

	var body bytes.Buffer
	if err := r.templates.ExecuteTemplate(&body, page.Template, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(page.Status)
	_, _ = w.Write(body.Bytes())
}
