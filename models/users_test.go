package models

import (
	"fmt"
	"testing"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	dbname   = "lenslocked_test"
	password = ""
)

func setup(t *testing.T) UserService {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	t.Log("psqlInfo:", psqlInfo)
	ug, err := NewUserGorm(psqlInfo)
	if err != nil {
		t.Fatal(err)
	}
	ug.LogMode(true)
	ug.DestructiveReset()
	return UserService(ug)
}

func teardown(t *testing.T, us UserService) {
	us.Close()
}

func TestUsers(t *testing.T) {
	us := setup(t)
	userData := User{
		Name:     "Test User",
		Password: "testpassword",
		Email:    "testuser@example.com",
	}

	t.Log("Should create a user", userData)

	if err := us.Create(&userData); err != nil {
		t.Fatalf("\tFailed to create a user for data (%v), error: %v",
			userData, err)
	}
	t.Log("\tCreated a user successfully ")

	t.Log("Should be able to sign a existing user successfully")
	if user, err := us.Authenticate("testuser@example.com", "testpassword"); err != nil {
		t.Error("\tFailed to authenticate  user", userData, err)
	} else {
		t.Log("\tUser authenticated successfully", user)
	}

	teardown(t, us)
}
