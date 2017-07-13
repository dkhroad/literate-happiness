package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

func NewServices(connectionInfo string) (*Services, error) {
	db, err := newGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		db:      db,
	}, nil
}

type Services struct {
	User    UserService
	Gallery GalleryService
	db      *gorm.DB
}

func newGorm(connectionInfo string) (*gorm.DB, error) {

	db, err := gorm.Open("postgres", connectionInfo)
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
	for _, model := range models {
		if err := svcs.db.DropTableIfExists(model).Error; err != nil {
			log.Fatal(err)
		}
	}
	svcs.AutoMigrate(models...)
}

func (svcs *Services) AutoMigrate(models ...interface{}) {
	// if err := svcs.db.AutoMigrate(models).Error; err != nil {
	if err := svcs.db.AutoMigrate(models...).Error; err != nil {
		log.Fatal(err)
	}
}

func (svcs *Services) Close() error {
	return svcs.db.Close()
}
