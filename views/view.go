package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

var (

	// LayoutDir holds the layouts root directory
	LayoutDir = "views/layouts/"
	// TemplateExt represents templates extension found on layout directory
	TemplateExt = ".gohtml"
	TemplateDir = "views/"
)

// NewView creates a new View
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)

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

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, nil)
}

// Render is used to render the view with the predefined Layout.
func (v *View) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	switch data.(type) {
	case Data:
		// do nothing
	default:
		data = Data{
			Yield: data,
		}
	}

	var buf bytes.Buffer

	if err := v.Template.ExecuteTemplate(&buf, v.Layout, data); err != nil {
		http.Error(w, "Something went wrong. If the problem persists, please email support@lenslocked.com",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)

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

// addTemplatePath takes in a slice of strings
// representing file paths for templates, and it prepends
// the TemplateDir directory to each string in the slice
//
// Eg the input {"home"} would result in the output
// {"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt takes in a slice of strings
// representing file paths for templates, and it appends
// the TemplateExt extension to each string in the slice
//
// Eg the input {"home"} would result in the output
// {"home.gohtml"} if TemplateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
