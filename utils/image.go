package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var maxSize int64 = 1024 * 1024 * 10

func UploadImageHandler(c *gin.Context, path string) (*string, error) {
	file, err := c.FormFile("image")

	if err != nil {
		return nil, ErrUploadImage
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !isAllowedExt(ext) {
		return nil, ErrUploadImageExt
	}

	if file.Size > maxSize {
		return nil, ErrUploadImageSize
	}

	savedFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(path, savedFilename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return nil, ErrSaveImage
	}

	return &savedFilename, nil
}

func isAllowedExt(ext string) bool {
	allowedExts := []string{".jpg", ".jpeg", ".png", ".svg"}

	for _, v := range allowedExts {
		if v == ext {
			return true
		}
	}

	return false
}

func ServeImage(c *gin.Context) {
	path := c.Param("path")
	fullPath := filepath.Join("./uploads", path)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		ErrorResponse(c, http.StatusNotFound, "Image not found")
		return
	}

	c.File(fullPath)
}

func UploadMultipleImageHandler(c *gin.Context, path string) ([]string, error) {
	form, err := c.MultipartForm()

	if err != nil {
		return nil, ErrUploadImage
	}

	files := form.File["images"]

	if len(files) == 0 {
		return nil, ErrEmptyUpload
	}

	var uploadedImages []string
	for _, file := range files {
		if file.Size > maxSize {
			return nil, ErrUploadImageSize
		}

		ext := filepath.Ext(file.Filename)
		if !isAllowedExt(ext) {
			return nil, ErrUploadImageExt
		}

		savedFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}

		fullPath := filepath.Join(path, savedFilename)

		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			return nil, ErrSaveImage
		}

		uploadedImages = append(uploadedImages, savedFilename)
	}

	return uploadedImages, nil
}
