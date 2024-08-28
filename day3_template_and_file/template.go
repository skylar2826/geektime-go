package __template_and_file

import (
	"bytes"
	"context"
	"text/template"
)

type templateEngine interface {
	Render(ctx context.Context, tagName string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tagName string, data any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(buffer, tagName, data)
	return buffer.Bytes(), err
}
