package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewServices(dialect string, connectionInfo string) (*Services, error) {
	db, err := newGorm(dialect, connectionInfo)
	if err != nil {
		return nil, err
	}
	// TODO: config this
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		Image:   NewImageService(),
		db:      db,
	}, nil
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
