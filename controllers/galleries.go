package controllers

import (
	"github.com/samueldaviddelacruz/lenslocked.com/models"
	"github.com/samueldaviddelacruz/lenslocked.com/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}
