package models

import (
	"github.com/jinzhu/gorm"
)

// Gallery is our image container resource that visitors
// see.
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

type GalleryService interface {
	GalleryDB
}

type galleryService struct {
	GalleryDB
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			&galleryGorm{db},
		},
	}
}

type GalleryDB interface {
	Create(gallery *Gallery) error
}

type galleryValidator struct {
	GalleryDB
}

var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

// Create will create the provided gallery and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
