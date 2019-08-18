package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/samueldaviddelacruz/lenslocked.com/controllers"
	"github.com/samueldaviddelacruz/lenslocked.com/middleware"
	"github.com/samueldaviddelacruz/lenslocked.com/models"
	"github.com/samueldaviddelacruz/lenslocked.com/rand"
)

func main() {
	boolPtr := flag.Bool("prod", false,
		"Provide this flag in production. This ensures that a config.json file is provided before the application starts")
	flag.Parse()
	appCfg := LoadConfig(*boolPtr)
	postgresConfig := appCfg.Database

	services, err := models.NewServices(
		models.WithGorm(
			postgresConfig.Dialect(),
			postgresConfig.ConnectionInfo()),
		models.WithLogMode(!appCfg.IsProd()),
		models.WithUser(appCfg.Pepper, appCfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
	)
	must(err)

	defer services.Close()
	//must(services.DestructiveReset())
	must(services.AutoMigrate())
	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)

	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	bytes, err := rand.Bytes(32)
	must(err)
	csrfMw := csrf.Protect(bytes, csrf.Secure(appCfg.IsProd()))

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	//assets
	assetsHandler := http.FileServer(http.Dir("./assets"))
	assetsHandler = http.StripPrefix("/assets/", assetsHandler)
	r.PathPrefix("/assets/").Handler(assetsHandler)
	// Image routes
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// Gallery routes

	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")

	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")

	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)

	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")

	// POST /galleries/:id/images/:filename/delete
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	fmt.Printf("Starting the server on port :%d\n", appCfg.Port)

	http.ListenAndServe(fmt.Sprintf(":%d", appCfg.Port), csrfMw(userMw.Apply(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
