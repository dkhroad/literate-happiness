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
		LoginView:   views.NewView("bootstrap", "users/login"),
		UserService: us,
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	loginForm := struct {
		Email    string `schema:"email"`
		Password string `schema:"password"`
	}{}

	if err := parseForm(r, &loginForm); err != nil {
		panic(err)
	}

	user, err := u.UserService.Authenticate(loginForm.Email, loginForm.Password)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, user)
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	form := SignupForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.UserService.Create(&user); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, user)
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	models.UserService
}
