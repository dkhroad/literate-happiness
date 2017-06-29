package controllers

import (
	"fmt"
	"lenslocked.com/models"
	"lenslocked.com/views"
	"net/http"
)

func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:     views.NewView("bootstrap", "users/new"),
		UserService: us,
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
	user := models.User{
		Name:  form.Name,
		Email: form.Email,
	}

	if err := u.UserService.Create(&user); err != nil {
		panic(err)
	}
	fmt.Println(form)
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type Users struct {
	NewView *views.View
	models.UserService
}
