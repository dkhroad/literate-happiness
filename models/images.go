package models

import (
	"fmt"
	"io"
	"os"
)

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct {
}

func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	galleryPath, err := is.createGalleryPath(galleryID)
	if err != nil {
		return err
	}
	imageFile := galleryPath + filename
	defer r.Close()

	var fd *os.File
	if fd, err = os.OpenFile(imageFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
		return err
	}
	defer fd.Close()
	if _, err = io.Copy(fd, r); err != nil {
		return err
	}
	return nil
}

func (is *imageService) createGalleryPath(galleryID uint) (string, error) {
	galleryPath := fmt.Sprintf("galleries/%v/images/", galleryID)
	if err := os.MkdirAll(galleryPath, 0755); err != nil {
		return "", err
	}
	return galleryPath, nil
}
