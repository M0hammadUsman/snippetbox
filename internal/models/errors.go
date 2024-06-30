package models

import "errors"

var ErrNoRows = errors.New("models: no matching records found")
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrDuplicateEmail = errors.New("models: duplicate email")
