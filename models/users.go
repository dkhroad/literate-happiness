package models

import (
	"fmt"
	"log"
	"regexp"

	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"lenslocked.com/hash"
	"lenslocked.com/rand"
)

const (
	pepperHash    = "doormat-wrangle-scam-gating-shelve"
	hmacSecretKey = "hmac-secret-key"
)

type UserService interface {
	UserDB
	Authenticate(string, string) (*User, error)
}

// compile time check to ensure UserService interface satisfies all interfaces
// defined by userService
var _ UserService = &userService{}

// UserDB is used to interact with a persistent store
type UserDB interface {
	// single user queries
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRememberTokenHash(token string) (*User, error)

	// CRUD methods
	Create(user *User) error
	Update(user *User) error
	UpdateAttributes(user *User) error
	Delete(id uint) error
}

func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{DB: db}
	uv := newUserValidator(ug, hash.NewHMAC(hmacSecretKey))
	return &userService{
		UserDB: uv,
	}
}

type userService struct {
	UserDB
}

func (us *userService) Authenticate(userEmail string, userPassword string) (*User, error) {
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

var _ UserDB = &userGorm{}

type userGorm struct {
	*gorm.DB
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	return ug.byQuery(ug.Where("id = ?", id))
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	log.Println("Looking for user with email: ", email)
	return ug.byQuery(ug.Where("email = ?", email))
}

func (ug *userGorm) ByRememberTokenHash(tokenHash string) (*User, error) {
	return ug.byQuery(ug.Where("remember_token_hash = ?", tokenHash))
}

func (ug *userGorm) Create(user *User) error {
	return ug.DB.Create(user).Error
}

func (ug *userGorm) UpdateAttributes(user *User) error {
	return ug.DB.Model(user).Updates(*user).Error
}

func (ug *userGorm) Update(user *User) error {
	return ug.DB.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	user := &User{Model: gorm.Model{ID: id}}
	return ug.DB.Delete(user).Error
}

func (ug *userGorm) Close() error {
	return ug.DB.Close()
}

func (ug *userGorm) byQuery(query *gorm.DB) (*User, error) {
	u := User{}
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

type User struct {
	gorm.Model
	Name              string
	Email             string `gorm:"not null;unique_index"`
	PasswordHash      string `gorm:"not null"`
	Password          string `gorm:"-"`
	RememberToken     string `gorm:"-"`
	RememberTokenHash string `gorm:"not null;unique_index"`
}

// validations and normalizations

type userValidator struct {
	UserDB
	hmac        hash.HMAC
	emailRegExp *regexp.Regexp
}

func newUserValidator(db UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:      db,
		hmac:        hmac,
		emailRegExp: regexp.MustCompile(`^[a-zA-Z%._-]+@[.a-zA-Z%_-]+\.[a-zA-Z]{2,16}$`),
	}
}
func (uv *userValidator) bcryptPassword(user *User) error {
	passwd := []byte(user.Password + pepperHash)
	phash, err := bcrypt.GenerateFromPassword(passwd, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(phash)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRememberToken(user *User) error {
	if user.RememberToken == "" {
		return nil
	}
	user.RememberTokenHash = uv.hmac.Hash(user.RememberToken)
	user.RememberToken = ""
	return nil
}

func (uv *userValidator) setRememberTokenIfUnset(user *User) error {
	if user.RememberToken != "" {
		return nil
	}

	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.RememberToken = token
	return nil
}

func (uv *userValidator) idGreaterThanN(n uint) func(*User) error {
	return func(user *User) error {
		if user.ID <= n {
			return ErrInvalidId
		}
		return nil
	}
}

func (uv *userValidator) normalizeEmail(user *User) error {
	if user.Email == "" {
		return nil
	}
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) validateEmailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if uv.emailRegExp.MatchString(user.Email) == false {
		return ErrInvalidEmail
	}
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err != nil { // didn't find the user
		return nil
	}

	// is it the same user?
	if user.ID == existing.ID {
		return nil
	}
	return ErrEmailNotAvail
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) Create(user *User) error {

	err := runUserValidateFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.validateEmailFormat,
		uv.emailIsAvail,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberTokenIfUnset,
		uv.hmacRememberToken,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	err := runUserValidateFuncs(user,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.validateEmailFormat,
		uv.emailIsAvail,
		uv.bcryptPassword,
		uv.passwordMinLength,
		uv.passwordHashRequired,
		uv.hmacRememberToken,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) UpdateAttributes(user *User) error {
	err := runUserValidateFuncs(user,
		uv.hmacRememberToken,
		uv.normalizeEmail,
		uv.validateEmailFormat,
		uv.emailIsAvail,
	)
	if err != nil {
		return err
	}

	return uv.UserDB.UpdateAttributes(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValidateFuncs(&user, uv.idGreaterThanN(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// ByRememberTokenHash will hash the given token and call the
// ByRememberTokenHash on subsequent DB layer
func (uv *userValidator) ByRememberTokenHash(token string) (*User, error) {
	user := User{
		RememberToken: token,
	}
	err := runUserValidateFuncs(&user,
		uv.hmacRememberToken,
	)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByRememberTokenHash(user.RememberTokenHash)
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValidateFuncs(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(user.Email)
}

type userValidateFunc func(*User) error

func runUserValidateFuncs(user *User, fns ...userValidateFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}

	return nil
}
