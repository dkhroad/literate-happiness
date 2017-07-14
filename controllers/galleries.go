package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

const (
	ShowGallery  = "show"
	IndexGallery = "index"
)

func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		gs:        gs,
		New:       views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		router:    r,
	}
}

type Galleries struct {
	New       *views.View
	IndexView *views.View
	ShowView  *views.View
	EditView  *views.View
	gs        models.GalleryService
	router    *mux.Router
}

type galleryForm struct {
	Title string `schema:"title"`
}

func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user, ok := context.User(r.Context())
	if !ok {
		log.Println("No signed in user found in the request context. This shouldn't  have happened." +
			"Redirecting to login page")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		return
	}

	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	form := galleryForm{}
	vd := views.Data{}
	if err := parseForm(r, &form); err != nil {
		vd.AddAlert(err)
		g.New.Render(w, r, vd)
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
		g.New.Render(w, r, vd)
		return
	}

	url, err := g.router.Get(IndexGallery).URL()
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
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var vd views.Data
	vd.Yield = gallery

	g.ShowView.Render(w, r, vd)
}

func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	form := galleryForm{}
	vd := views.Data{}
	if err := parseForm(r, &form); err != nil {
		vd.AddAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	gallery.Title = form.Title
	if err := g.gs.Update(gallery); err != nil {
		vd.AddAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	vd.Alert = views.AlertSuccess("Gallery updates successfully")
	vd.Yield = gallery
	g.EditView.Render(w, r, vd)

}

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, r, vd)
}

func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	if err = g.gs.Delete(gallery.ID); err != nil {
		var vd views.Data
		vd.AddAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, r, vd)
	}

	url, err := g.router.Get(IndexGallery).URL()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	log.Println("redirecting to ", url.Path)
	http.Redirect(w, r, url.Path, http.StatusFound)
}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		err := errors.New("invalid or missing gallery id")
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid or missing gallery id", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		log.Println(err)
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return nil, err
	}
	return gallery, nil
}
