package model

import "errors"

var ErrEmptyUUID error = errors.New("generated uuid is empty")
