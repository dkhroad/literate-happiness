package views

import (
	"bytes"
	"html/template"
	"io"
	"log"
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
		log.Panic(err)
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
	v.Render(w, nil)
}

// Render is used to render a view with a predefined layout.
func (v *View) Render(w http.ResponseWriter, data interface{}) {
	switch data.(type) {
	case Data:
		// do nothing
	default:
		data = Data{
			Yield: data,
		}
	}
	// v.RenderError(w, data)

	w.Header().Set("Content-Type", "text/html")
	v.executeTemplateAndLogError(w, data)
}

func (v *View) RenderError(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html")
	v.executeTemplateAndLogError(w, data)
}

func (v *View) executeTemplateAndLogError(w http.ResponseWriter, data interface{}) {
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, data); err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if _, err := io.Copy(w, &buf); err != nil {
		log.Println("Failed to execute template", v.Layout, err)
	}
	return
}

func addTemplateDirAndExt(files []string) {
	for i, file := range files {
		files[i] = filepath.Join(TemplateDir, file) + TemplateExt
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(filepath.Join(LayoutDir, "/*"+TemplateExt))
	if err != nil {
		log.Panic(err)
	}
	return files
}
