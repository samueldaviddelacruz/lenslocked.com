package views

import (
	"html/template"
	"path/filepath"
)

var (
	LayoutDir   = "views/layouts/"
	TemplateExt = ".gohtml"
)

// NewView creates a new View
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

// View Represents a view template
type View struct {
	Template *template.Template
	Layout   string
}

// layoutfiles returns a slice of strings representing
// the layout files used in our application.
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)

	if err != nil {
		panic(err)
	}
	return files
}
