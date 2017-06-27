package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	TemplateDir = "views"
	LayoutDir   = "views/layouts"
	TemplateExt = ".tmpl"
)

func NewView(layout string, files ...string) *View {
	addTemplateDirAndExt(files)
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

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := v.Template.ExecuteTemplate(w, v.Layout, nil); err != nil {
		panic(err)
	}
}
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func addTemplateDirAndExt(files []string) {
	for i, file := range files {
		files[i] = filepath.Join(TemplateDir, file) + TemplateExt
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(filepath.Join(LayoutDir, "/*"+TemplateExt))
	if err != nil {
		panic(err)
	}
	return files
}
