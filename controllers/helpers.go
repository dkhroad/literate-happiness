package controllers

import (
	"github.com/gorilla/schema"
	"net/http"
)

func parseForm(r *http.Request, dst interface{}) error {
	var err error
	if err = r.ParseForm(); err != nil {
		return err
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(dst, r.PostForm)
	return err
}
