package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const imagePath = "images/galleries"

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(id uint) ([]Image, error)
	DeleteImage(img Image)
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct {
}

func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	galleryPath := imagePathWithGallery(galleryID)
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

func (is *imageService) ByGalleryID(id uint) ([]Image, error) {
	galleryPath := imagePathWithGallery(id)
	imgs, err := filepath.Glob(galleryPath + "*")
	if err != nil {
		return nil, err
	}
	images := make([]Image, len(imgs))
	for i, img := range imgs {
		fn := strings.Replace(img, galleryPath, "", 1)
		images[i] = Image{fn, id}
	}
	return images, nil
}

func (is *imageService) DeleteImage(img Image) {
	os.RemoveAll(img.RelativePath())
}

type Image struct {
	Filename  string
	GalleryID uint
}

func (im *Image) Path() string {
	return "/" + im.RelativePath()
}

func (im *Image) RelativePath() string {
	return imagePathWithGallery(im.GalleryID) + im.Filename
}

func imagePathWithGallery(id uint) string {
	return fmt.Sprintf("%v/%v/", imagePath, id)
}
