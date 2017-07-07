package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

func NewUsers(us *models.UserService) *Users {
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
		switch err {
		case models.ErrNotFound:
			fmt.Fprintf(w, "Invalid email address")
		case models.ErrInvalidPassword:
			fmt.Fprintf(w, "Invalid password")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := u.signInUser(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := u.signInUser(w, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *Users) signInUser(w http.ResponseWriter, user *models.User) error {
	fmt.Println(user)
	if user.RememberToken == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}

		// will generate and save the remember token hash
		if err := u.UserService.UpdateAttributes(user, models.User{
			RememberToken: token}); err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.RememberTokenHash,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	*models.UserService
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	user, err := u.ByRememberTokenHash(cookie.Value)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Fprintln(w, "remember token: ", cookie.Value)
	fmt.Fprintf(w, "%+v", user)
}
