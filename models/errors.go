package models

const (
	ErrNotFound         modelError = "Resource not found"
	ErrInvalidPassword  modelError = "Invalid password"
	ErrInvalidId        modelError = "Invalid used id"
	ErrEmailRequired    modelError = "Email is required"
	ErrInvalidEmail     modelError = "Email is not valid"
	ErrEmailNotAvail    modelError = "Email address is not available"
	ErrPasswordRequired modelError = "Password is required"
	ErrPasswordTooShort modelError = "Password is too short"
)

type modelError string

func (me modelError) Error() string {
	return string(me)
}

func (me modelError) Public() string {
	return string(me)
}
