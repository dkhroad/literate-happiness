package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
	"lenslocked.com/models"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	dbname   = "lenslocked_dev"
	password = ""
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	ug, err := models.NewUserGorm(psqlInfo)
	if err != nil {
		panic(err)
	}
	ug.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(ug)
	galleriesC := controllers.NewGalleries()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/galleries", galleriesC.New).Methods("GET")
	http.ListenAndServe(":3000", r)
}
