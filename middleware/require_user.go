package middleware

import (
	"log"
	"net/http"
	"strings"

	"lenslocked.com/context"
	"lenslocked.com/models"
)

type RequireUser struct {
	User
}

type User struct {
	models.UserService
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// skip this check for static assets
		path := r.URL.Path
		if strings.HasPrefix(path, "/images/") ||
			strings.HasPrefix(path, "/assets/") {
			next(w, r)
			return
		}
		// if user is logged in
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			log.Println("Remember token not found")
			next(w, r)
			return
		}
		user, err := mw.UserService.ByRememberTokenHash(cookie.Value)
		if user == nil {
			log.Println("Unable to found a user with the given remember token", cookie.Value)
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		log.Printf("Add user %v to %v\n", user.Name, r.URL)
		next(w, r)
	})
}
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := context.User(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
