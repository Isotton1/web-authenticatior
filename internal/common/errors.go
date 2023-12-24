package common

import "errors"

var ErrNoUserFound error = errors.New("no user found with the provided username")
