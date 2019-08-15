package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/samueldaviddelacruz/lenslocked.com/context"
	"github.com/samueldaviddelacruz/lenslocked.com/models"
	"github.com/samueldaviddelacruz/lenslocked.com/views"
)

const (
	ShowGallery = "show_gallery"
	EditGallery = "edit_gallery"

	maxMultiPartMem = 1 << 20 // 1 megabyte
)

func NewGalleries(gs models.GalleryService, is models.ImageService, router *mux.Router) *Galleries {
	return &Galleries{
		New:       views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		gs:        gs,
		is:        is,
		router:    router,
	}
}

type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	is        models.ImageService
	router    *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)

	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
	//fmt.Fprintln(w, galleries)

}

// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery

	g.ShowView.Render(w, r, vd)

}

// GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery

	g.EditView.Render(w, r, vd)

}

// POST /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery

	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery successfully updated!",
	}
	g.EditView.Render(w, r, vd)
}

// POST /galleries/:id/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	err = g.gs.Delete(gallery.ID)
	if err != nil {
		var vd views.Data
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)

}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}
	url, err := g.router.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))

	if err != nil {
		// TODO: make this go to index page
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// POST /galleries/:id/images
func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	// TODO: Parse a multipart form
	err = r.ParseMultipartForm(maxMultiPartMem)
	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	files := r.MultipartForm.File["images"]

	for _, f := range files {
		// Open the uploaded file
		file, err := f.Open()
		if err != nil {
			log.Println(err)
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer file.Close()
		err = g.is.Create(gallery.ID, file, f.Filename)
		if err != nil {
			log.Println(err)
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
	}
	url, err := g.router.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// POST /galleries/:id/images/:filename/delete
func (g *Galleries) ImageDelete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	filename := mux.Vars(r)["filename"]

	i := models.Image{
		Filename:  filename,
		GalleryID: gallery.ID,
	}

	err = g.is.Delete(&i)
	if err != nil {
		var vd views.Data
		vd.Yield = gallery

		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	url, err := g.router.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))

	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)

}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)

		default:
			http.Error(w, "Oops, something went wrong.", http.StatusInternalServerError)

		}
		return nil, err

	}
	images, _ := g.is.ByGalleryID(gallery.ID)
	gallery.Images = images

	return gallery, nil
}
