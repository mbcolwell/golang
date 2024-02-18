package cyoa

import (
	"embed"
	"html/template"
)

//go:embed *.json
var jsonFiles embed.FS

//go:embed *.html
var htmlFiles embed.FS

var tpl *template.Template

func init() {
	// If this doesn't run, the program should fall over before creating a server
	defaultHandlerTmpl, err := htmlFiles.ReadFile("story.html")
	if err != nil {
		panic(err)
	}
	tpl = template.Must(template.New("").Parse(string(defaultHandlerTmpl)))
}
