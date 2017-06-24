package views

import "html/template"

func NewView(layout string, files ...string) *View {
	files = append(files,
		"views/layouts/footer.tmpl",
		"views/layouts/bootstrap.tmpl",
		"views/layouts/navbar.tmpl",
	)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}
