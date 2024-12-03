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

func UploadImageHandler(c *gin.Context, path string) (*string, error) {
	var maxSize int64 = 1024 * 1024 * 10
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
