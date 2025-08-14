package models

import (
	"errors"
	"fmt"
)

var ErrNoInCache = errors.New("no data in cache")

type Error struct {
	Message string `json:"message"`
}

func Err(msg string) Error {
	return Error{
		Message: msg,
	}
}

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
