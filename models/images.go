package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(id uint) ([]string, error)
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct {
}

func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	galleryPath := is.getGalleryImagePath(galleryID)
	var err error
	if err = os.MkdirAll(galleryPath, 0755); err != nil {
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

func (is *imageService) ByGalleryID(id uint) ([]string, error) {
	galleryPath := is.getGalleryImagePath(id)
	return filepath.Glob(galleryPath + "*")
}

func (is *imageService) getGalleryImagePath(galleryID uint) string {
	return fmt.Sprintf("galleries/%v/images/", galleryID)
}
