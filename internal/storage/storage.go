package storage

import "errors"

var (
	ResourceNotFound      = errors.New("the required resource was not found")
	ResourceAlreadyExists = errors.New("this url already exists")
)
