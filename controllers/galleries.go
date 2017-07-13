package controllers

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		gs:  gs,
		New: views.NewView("bootstrap", "galleries/new"),
	}
}

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

type galleryForm struct {
	Title string `schema:"title"`
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	form := galleryForm{}
	vd := views.Data{}
	if err := parseForm(r, &form); err != nil {
		vd.AddAlert(err)
		g.New.Render(w, vd)
		return
	}

	user, ok := context.User(r.Context())
	if !ok {
		log.Println("No signed in user found in the request context. This shouldn't  have happened." +
			"Redirecting to login page")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.AddAlert(err)
		g.New.Render(w, vd)
		return
	}

	fmt.Fprintln(w, gallery)
}
