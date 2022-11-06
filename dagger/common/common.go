package common

import "errors"

// ErrMissingRequiredArgument is returned when a required argument is missing.
var ErrMissingRequiredArgument = errors.New("missing required argument")
