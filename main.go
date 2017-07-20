package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
	"lenslocked.com/middleware"
	"lenslocked.com/models"
	"lenslocked.com/rand"
)

func init() {
	log.SetPrefix("LENS: ")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	dbCfg := DefaultPostgresConfig()
	cfg := DefaultConfig()
	log.Println("config: ", cfg)
	log.Println("DBConfig: ", dbCfg.ConnectionInfo())

	svcs, err := models.NewServices(dbCfg.Dialect(), dbCfg.ConnectionInfo())
	if err != nil {
		log.Panic(err)
	}
	defer svcs.Close()
	// svcs.DestructiveReset()
	svcs.AutoMigrate()

	r := mux.NewRouter()

	// CSRF protection
	randomBytes, err1 := rand.Bytes(32)
	if err1 != nil {
		log.Panic(err1)
	}

	csrfMW := csrf.Protect(randomBytes, csrf.Secure(!cfg.isProd()))

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(svcs.User)
	galleriesC := controllers.NewGalleries(svcs.Gallery, svcs.Image, r)

	assetHandler := http.FileServer(http.Dir("./assets"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets").Handler(assetHandler)

	imageHandler := http.FileServer(http.Dir("./images"))
	r.PathPrefix("/images").Handler(http.StripPrefix("/images/", imageHandler))

	userMw := middleware.User{UserService: svcs.User}
	requireUserMw := middleware.RequireUser{User: userMw}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.HandleFunc("/galleries",
		requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET").Name(controllers.IndexGallery)

	r.Handle("/galleries/new",
		requireUserMw.Apply(galleriesC.New)).Methods("GET")

	r.HandleFunc("/galleries",
		requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/edit",
		requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)

	r.HandleFunc("/galleries/{id:[0-9]+}/update",
		requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}",
		requireUserMw.ApplyFn(galleriesC.Show)).Methods("GET").Name(controllers.ShowGallery)

	r.HandleFunc("/galleries/{id:[0-9]+}/images",
		requireUserMw.ApplyFn(galleriesC.UploadImages)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		requireUserMw.ApplyFn(galleriesC.DeleteImage)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/delete",
		requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMW(userMw.Apply(r)))
}
