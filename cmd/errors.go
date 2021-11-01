package cmd

import (
	"errors"
)

var errNotEnoughArguments = errors.New("not enough arguments")

var errTooManyArguments = errors.New("too many arguments")
