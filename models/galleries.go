package models

import "github.com/jinzhu/gorm"

type Gallery struct {
	gorm.Model
	Title  string  `gorm:"not null;unique_index"`
	UserID uint    `gorm:"not null;index"`
	Images []Image `gorm:"-"`
}

func (g *Gallery) ImageSplitN(n int) [][]Image {
	ret := make([][]Image, n)
	for i := 0; i < n; i++ {
		ret[i] = make([]Image, 0)
	}

	for i, img := range g.Images {
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)
	}
	return ret
}

type GalleryDB interface {
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
	ByID(id uint) (*Gallery, error)
	ByUserID(id uint) ([]Gallery, error)
}

type GalleryService interface {
	GalleryDB
}

type galleryService struct {
	GalleryDB
}

var _ GalleryService = &galleryService{}
var _ GalleryDB = &galleryGorm{}
var _ GalleryDB = &galleryValidator{}

func NewGalleryService(db *gorm.DB) GalleryService {
	gg := &galleryGorm{db: db}
	gv := &galleryValidator{GalleryDB: gg}
	return &galleryService{GalleryDB: gv}
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) Delete(id uint) error {
	gallery := &Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(gallery).Error
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	return gg.byQuery(gg.db.Where("id = ?", id))
}

func (gg *galleryGorm) ByUserID(id uint) ([]Gallery, error) {
	var galleries []Gallery
	err := gg.db.Where("user_id = ?", id).Find(&galleries).Error
	return galleries, err
}

func (gg *galleryGorm) byQuery(query *gorm.DB) (*Gallery, error) {
	u := Gallery{}
	err := query.First(&u).Error
	switch err {
	case nil:
		return &u, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

//
// galleryValidator
//
type galleryValidator struct {
	GalleryDB
}

func (gv *galleryValidator) titleRequired(gallery *Gallery) error {
	if gallery.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (gv *galleryValidator) userIDRequired(gallery *Gallery) error {
	if gallery.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) idGreaterThanN(n uint) func(*Gallery) error {
	return func(gallery *Gallery) error {
		if gallery.ID <= n {
			return ErrInvalidId
		}
		return nil
	}
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	// TODO add uniqueness validator
	err := runGalleryValidatorFuncs(gallery,
		gv.titleRequired,
		gv.userIDRequired,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

func (gv *galleryValidator) Update(gallery *Gallery) error {
	err := runGalleryValidatorFuncs(gallery,
		gv.titleRequired,
		gv.userIDRequired,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

func (gv *galleryValidator) Delete(id uint) error {
	var gallery Gallery
	gallery.ID = id
	err := runGalleryValidatorFuncs(&gallery, gv.idGreaterThanN(0))
	if err != nil {
		return err
	}
	return gv.GalleryDB.Delete(id)
}

type galleryValidatorFuncs func(gallery *Gallery) error

func runGalleryValidatorFuncs(gallery *Gallery, fns ...galleryValidatorFuncs) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}
