package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

const (
	ShowGallery = "show"
)

func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		gs:       gs,
		New:      views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
		router:   r,
	}
}

type Galleries struct {
	New      *views.View
	ShowView *views.View
	gs       models.GalleryService
	router   *mux.Router
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

	// fmt.Fprintln(w, gallery)
	url, err := g.router.Get(ShowGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	log.Println("redirecting to ", url.Path)
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		http.Error(w, "invalid or missing gallery id", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid or missing gallery id", http.StatusNotFound)
		return
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		log.Println(err)
		http.Error(w, "Gallery not found", http.StatusNotFound)
	}
	var vd views.Data
	vd.Yield = gallery

	g.ShowView.Render(w, vd)
}
