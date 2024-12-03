package utils

import "errors"

var ErrAlreadyExists = errors.New("already exists")
var ErrNotFound = errors.New("not found")
var ErrUploadImage = errors.New("failed to upload image")
var ErrUploadImageExt = errors.New("image extension not allowed")
var ErrUploadImageSize = errors.New("image size too large")
var ErrSaveImage = errors.New("failed to save image")
