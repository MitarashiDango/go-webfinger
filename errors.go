package webfinger

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidResponse  = errors.New("invalid response")
	ErrResourceNotFound = errors.New("resource not found")
)

type Error struct {
	Err error
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	if e.Err == nil {
		return "webfinger error"
	}

	return fmt.Sprintf("webfinger error: %s", e.Err)
}

type WebFingerResponseStatusError struct {
	URL        string
	StatusCode int
	Status     string
}

func (e *WebFingerResponseStatusError) Error() string {
	return fmt.Sprintf("webfinger response status error: %d %s", e.StatusCode, e.Status)
}

type UnsupportedContentTypeError struct {
	ContentType string
}

func (e *UnsupportedContentTypeError) Error() string {
	return fmt.Sprintf("unsupported content type error: %s", e.ContentType)
}
