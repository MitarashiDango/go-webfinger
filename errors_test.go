package webfinger_test

import (
	"errors"
	"testing"

	webfinger "github.com/MitarashiDango/go-webfinger"
)

func Test_Error_Unwrap_001(t *testing.T) {
	e := &webfinger.Error{}
	if err := e.Unwrap(); err != nil {
		t.FailNow()
	}
}

func Test_Error_Unwrap_002(t *testing.T) {
	e := &webfinger.Error{
		Err: errors.New("test error"),
	}
	if err := e.Unwrap(); err == nil {
		t.FailNow()
	}
}

func Test_Error_Error_001(t *testing.T) {
	e := &webfinger.Error{}
	if err := e.Error(); err != "webfinger error" {
		t.FailNow()
	}
}

func Test_Error_Error_002(t *testing.T) {
	e := &webfinger.Error{
		Err: errors.New("test error"),
	}
	if err := e.Error(); err != "webfinger error: test error" {
		t.FailNow()
	}
}

func Test_WebFingerResponseStatusError_Error_001(t *testing.T) {
	e := &webfinger.WebFingerResponseStatusError{
		StatusCode: 100,
		Status:     "test status",
	}
	if err := e.Error(); err != "webfinger response status error: 100 test status" {
		t.FailNow()
	}
}

func Test_UnsupportedContentTypeError_Error_001(t *testing.T) {
	e := &webfinger.UnsupportedContentTypeError{
		ContentType: "test content type",
	}
	if err := e.Error(); err != "unsupported content type error: test content type" {
		t.FailNow()
	}
}
