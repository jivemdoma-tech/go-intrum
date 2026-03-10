package gointrum

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrEmptySubdomain      = errors.New("empty subdomain")
	ErrEmptyApiKey         = errors.New("empty api key")
	ErrEmptyParams         = errors.New("empty params")
	ErrEmptyRequiredParams = errors.New("missing required params")
)

var (
	ErrNothingFound = errors.New("nothing found")
)

func newErr(methodURL string, msg error) error {
	switch u, err := url.ParseRequestURI(methodURL); {
	case err != nil:
		return fmt.Errorf("failed to complete request: %w", msg)
	default:
		return fmt.Errorf("failed to complete request %s: %w", u.Path, msg)
	}
}

func validateRequestArgs(methodURL, subdomain, apiKey string) error {
	switch {
	case subdomain == "":
		return newErrEmptySubdomain(methodURL)
	case apiKey == "":
		return newErrEmptyApiKey(methodURL)
	default:
		return nil
	}
}

func newErrEmptySubdomain(methodURL string) error {
	return newErr(methodURL, ErrEmptySubdomain)
}

func newErrEmptyApiKey(methodURL string) error {
	return newErr(methodURL, ErrEmptyApiKey)
}

func newErrEmptyParams(methodURL string) error {
	return newErr(methodURL, ErrEmptyParams)
}

func newErrEmptyRequiredParams(methodURL string) error {
	return newErr(methodURL, ErrEmptyRequiredParams)
}
