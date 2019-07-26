package views

import "html/template"

// NewView creates a new View
func NewView(layout string, files ...string) *View {
	files = append(files,
		"views/layouts/navbar.gohtml",
		"views/layouts/bootstrap.gohtml",
		"views/layouts/footer.gohtml")
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
