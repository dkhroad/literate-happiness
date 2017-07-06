package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"lenslocked.com/hash"
)

var (
	ErrNotFound        = errors.New("models: resource not found")
	ErrInvalidPassword = errors.New("Invalid password")
	ErrInvalidId       = errors.New("Invalid used id")
)

const (
	pepperHash    = "doormat-wrangle-scam-gating-shelve"
	hmacSecretKey = "hmac-secret-key"
)

// UserDB is used to interact with a persistent store
type UserDB interface {
	// single user queries
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRememberTokenHash(token string) (*User, error)

	// CRUD methods
	Create(user *User) error
	Update(user *User) error
	UpdateAttributes(user *User, attrs User) error
	Delete(id uint) error

	// close the database connection to avoid resource leakage
	Close() error

	// migrations
	DestructiveReset()
	AutoMigrate()
}

func NewUserService(connectionInfo string) (*UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}

	return &UserService{
		UserDB: &UserValidator{
			UserDB: ug,
		},
	}, nil
}

type UserService struct {
	UserDB
}

type UserValidator struct {
	UserDB
}

func newUserGorm(connectionInfo string) (*UserGorm, error) {

	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserGorm{
		DB:   db,
		hmac: hash.NewHMAC(hmacSecretKey),
	}, nil
}

func (ug *UserGorm) ByID(id uint) (*User, error) {
	return ug.byQuery(ug.Where("id = ?", id))
}

func (ug *UserGorm) ByEmail(email string) (*User, error) {
	return ug.byQuery(ug.Where("email = ?", email))
}

func (ug *UserGorm) ByRememberTokenHash(tokenHash string) (*User, error) {
	return ug.byQuery(ug.Where("remember_token_hash = ?", tokenHash))
}

func (ug *UserGorm) Create(user *User) error {
	passwd := []byte(user.Password + pepperHash)
	phash, err := bcrypt.GenerateFromPassword(passwd, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(phash)
	user.Password = ""

	return ug.DB.Create(user).Error
}

func (ug *UserGorm) UpdateAttributes(user *User, attrs User) error {
	if attrs.RememberToken != "" {
		attrs.RememberTokenHash = ug.hmac.Hash(attrs.RememberToken)
	}
	return ug.DB.Model(user).Updates(attrs).Error
}

func (ug *UserGorm) Update(user *User) error {
	if user.RememberToken != "" {
		user.RememberTokenHash = ug.hmac.Hash(user.RememberToken)
	}
	return ug.DB.Save(user).Error
}

func (ug *UserGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidId
	}
	user := &User{Model: gorm.Model{ID: id}}
	return ug.DB.Delete(user).Error
}

func (ug *UserGorm) Close() error {
	return ug.DB.Close()
}

func (ug *UserGorm) byQuery(query *gorm.DB) (*User, error) {
	u := User{}
	err := query.First(&u).Error
	switch err {
	case nil:
		return &u, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		panic(err)
	}
}

func (ug *UserGorm) DestructiveReset() {
	ug.DropTableIfExists(&User{})
	ug.AutoMigrate()
}

func (ug *UserGorm) AutoMigrate() {
	ug.DB.AutoMigrate(&User{})
}

var _ UserDB = &UserGorm{}

type UserGorm struct {
	*gorm.DB
	hmac hash.HMAC
}

type User struct {
	gorm.Model
	Name              string
	Email             string `gorm:"not null;unique_index"`
	PasswordHash      string `gorm:"not null"`
	Password          string `gorm:"-"`
	RememberToken     string `gorm:"-"`
	RememberTokenHash string `gorm:"not null;unique_index"`
}

func (us *UserService) Authenticate(userEmail string, userPassword string) (*User, error) {
	foundUser, err := us.ByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(userPassword+pepperHash))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	fmt.Println(foundUser)
	return foundUser, err
}
