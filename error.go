package gointrum

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrNilParams           = errors.New("nil request params")
	ErrEmptyRequiredFields = errors.New("empty required fields")
)

func newErrNilParams(methodURL string) error {
	return newErr(methodURL, ErrNilParams)
}

func newErrEmptyRequiredFields(methodURL string) error {
	return newErr(methodURL, ErrEmptyRequiredFields)
}

func newErr(methodURL string, msg error) error {
	u, _ := url.ParseRequestURI(methodURL)
	return fmt.Errorf("failed to process request %s: %w", u.Path, msg)
}
