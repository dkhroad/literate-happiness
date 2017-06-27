package controllers

import (
	"fmt"
	"lenslocked.com/views"
	"net/http"
)

func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In Create users...")
	form := SignupForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	fmt.Println(form)
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type Users struct {
	NewView *views.View
}
