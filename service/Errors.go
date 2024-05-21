package service

import "errors"

// common errors

var ErrorStoreError = errors.New("store error")

var ErrorNotFound = errors.New("not found")

var ErrorValidation = errors.New("validation error")
