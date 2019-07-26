package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	// LayoutDir holds the layouts root directory
	LayoutDir = "views/layouts/"
	// TemplateExt represents templates extension found on layout directory
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

// Render is used to render the view with the predefined Layout.
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
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
