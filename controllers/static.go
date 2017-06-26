package controllers

import "lenslocked.com/views"

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "views/static/home.tmpl"),
		Contact: views.NewView("bootstrap", "views/static/contact.tmpl"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
