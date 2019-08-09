package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samueldaviddelacruz/lenslocked.com/controllers"
	"github.com/samueldaviddelacruz/lenslocked.com/middleware"
	"github.com/samueldaviddelacruz/lenslocked.com/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "lenslocked_dev"
)

func main() {
	psqlinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices(psqlinfo)
	must(err)

	defer services.Close()
	//must(services.DestructiveReset())
	must(services.AutoMigrate())
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)

	galleriesC := controllers.NewGalleries(services.Gallery)
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}
	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	r.HandleFunc("/cookieTest", usersC.CookieTest).Methods("GET")

	// Gallery routes

	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	fmt.Println("Starting the server on port :4000")
	http.ListenAndServe(":4000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
