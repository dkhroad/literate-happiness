package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func WithUserGorm(dialect string, connectionInfo string) func(*Services) error {
	return func(s *Services) error {
		db, err := newGorm(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithUser() func(*Services) error {
	return func(s *Services) error {
		s.User = NewUserService(s.db)
		return nil
	}
}

func WithLogMode(mode bool) func(*Services) error {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

func WithGallery() func(*Services) error {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

func WithImage() func(*Services) error {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

func NewServices(cfgs ...func(*Services) error) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

type Services struct {
	User    UserService
	Gallery GalleryService
	Image   ImageService
	db      *gorm.DB
}

func newGorm(dialect string, connectionInfo string) (*gorm.DB, error) {
	// TODO: config this
	db, err := gorm.Open(dialect, connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return db, nil
}

var models = []interface{}{
	&User{},
	&Gallery{},
}

func (svcs *Services) DestructiveReset() {
	log.Println("DestructiveReset", models)
	for _, model := range models {
		if err := svcs.db.DropTableIfExists(model).Error; err != nil {
			log.Fatal(err)
		}
	}
	svcs.AutoMigrate()
}

func (svcs *Services) AutoMigrate() {
	log.Println("auto migrating..", models)
	if err := svcs.db.AutoMigrate(models...).Error; err != nil {
		log.Fatal(err)
	}
}

func (svcs *Services) Close() error {
	return svcs.db.Close()
}
