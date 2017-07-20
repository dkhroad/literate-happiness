package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"

	"lenslocked.com/context"
)

var (
	TemplateDir = "views"
	LayoutDir   = "views/layouts"
	TemplateExt = ".tmpl"
)

func NewView(layout string, files ...string) *View {
	addTemplateDirAndExt(files)
	files = append(files, layoutFiles()...)

	// t, err := template.ParseFiles(files...)
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrfField is not Implemented")
		},
	}).ParseFiles(files...)

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
	v.Render(w, r, nil)
}

// Render is used to render a view with a predefined layout.
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	var vd Data
	switch d := data.(type) {
	case Data:
		// do nothing
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	vd.User, _ = context.User(r.Context())
	w.Header().Set("Content-Type", "text/html")
	v.executeTemplateAndLogError(tpl, w, vd)
}

func (v *View) RenderError(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html")
	v.executeTemplateAndLogError(v.Template, w, data)
}

func (v *View) executeTemplateAndLogError(tpl *template.Template, w http.ResponseWriter, data interface{}) {
	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, v.Layout, data); err != nil {
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
