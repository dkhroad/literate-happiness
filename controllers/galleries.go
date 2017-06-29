package controllers

import (
	"fmt"
	"net/http"
)

func NewGalleries() *Galleries {
	return &Galleries{}
}

func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is where you will create galleries resource")
}

type Galleries struct {
}
