package views

import "html/template"

// NewView creates a new View
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
	}
}

// View Represents a view template
type View struct {
	Template *template.Template
}
