package user

import (
	"embed"
	"html/template"
)

//go:embed templates/*.html templates/partials/*.html
var templateFS embed.FS

func Templates() *template.Template {
	return template.Must(template.ParseFS(templateFS, "templates/*.html", "templates/partials/*.html"))
}
