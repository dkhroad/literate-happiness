package views

import "html/template"

func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.tmpl")
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{t}
}

type View struct {
	Template *template.Template
}
