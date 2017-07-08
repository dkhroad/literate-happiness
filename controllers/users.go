package controllers

import (
	"fmt"
	"log"
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
		log.Fatal(err)
		vd := views.Data{
			Alert: views.AlertGeneric,
		}
		u.NewView.RenderError(w, vd)
	}
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	form := SignupForm{}
	vd := views.Data{}
	if err := parseForm(r, &form); err != nil {
		log.Fatal(err)
		vd.AddAlert(views.AlertGeneric)
		u.NewView.RenderError(w, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.UserService.Create(&user); err != nil {
		vd.AddAlert(views.AlertError(err.Error()))
		u.NewView.RenderError(w, vd)
		return
	}
	if err := u.signInUser(w, &user); err != nil {
		log.Fatal(err)
		vd.AddAlert(views.AlertWarning("Your account was created, but we were unable to log you in." +
			"Please try to login again."))
		u.NewView.RenderError(w, vd)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	loginForm := struct {
		Email    string `schema:"email"`
		Password string `schema:"password"`
	}{}

	vd := view.Data{}
	if err := parseForm(r, &loginForm); err != nil {
		log.Fatal(err)
		vd.AddAlert(views.AlertGeneric)
		u.LoginView.RenderError(w, vd)
		return
	}

	user, err := u.UserService.Authenticate(loginForm.Email, loginForm.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AddAlert(views.AlertError("Invalid email address"))
		case models.ErrInvalidPassword:
			vd.AddAlert(views.AlertError("Invalid password"))
		default:
			vd.AddAlert(views.AlertGeneric)
		}
		u.LoginView.RenderError(w, vd)
		return
	}
	if err := u.signInUser(w, user); err != nil {
		log.Fatal(err)
		vd.AddAlert(views.AlertGeneric)
		u.LoginView.RenderError(w, vd)
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
