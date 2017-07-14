package controllers

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	form := SignupForm{}
	vd := views.Data{}
	if err := parseForm(r, &form); err != nil {
		vd.AddAlert(err)
		u.NewView.RenderError(w, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		vd.AddAlert(err)
		u.NewView.RenderError(w, vd)
		return
	}
	if err := u.signInUser(w, &user); err != nil {

		log.Println("User account was created, but we were unable to log the user in. Redirecting to the login page", user, err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	loginForm := struct {
		Email    string `schema:"email"`
		Password string `schema:"password"`
	}{}

	vd := views.Data{}
	if err := parseForm(r, &loginForm); err != nil {
		vd.AddAlert(err)
		u.LoginView.RenderError(w, vd)
		return
	}

	user, err := u.us.Authenticate(loginForm.Email, loginForm.Password)
	if err != nil {
		vd.AddAlert(err)
		u.LoginView.RenderError(w, vd)
		return
	}

	if err := u.signInUser(w, user); err != nil {
		vd.AddAlert(err)
		u.LoginView.RenderError(w, vd)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *Users) signInUser(w http.ResponseWriter, user *models.User) error {
	if user.RememberToken == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.RememberToken = token
		// will generate and save the remember token hash
		if err := u.us.UpdateAttributes(&models.User{RememberToken: token}); err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.RememberToken,
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
	us        models.UserService
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	user, err := u.us.ByRememberTokenHash(cookie.Value)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Fprintln(w, "remember token: ", cookie.Value)
	fmt.Fprintf(w, "%+v", user)
}
