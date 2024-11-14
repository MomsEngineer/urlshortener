package internalerrors

import "errors"

var ErrDuplicate = errors.New("duplicate entry")
var ErrNoContent = errors.New("no content")
