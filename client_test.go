package webfinger_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	webfinger "github.com/MitarashiDango/go-webfinger"
)

func Test_Client_Do_JSONResponse(t *testing.T) {

	var baseURL string
	var host string
	renderWebFingerResponse := func() string {
		return `
		{
			"subject": "acct:test@` + host + `",
			"aliases": [
					"` + baseURL + `/@test",
					"` + baseURL + `/users/test"
			],
			"links": [
					{
							"rel": "self",
							"type": "application/activity+json",
							"href": "` + baseURL + `/users/test"
					}
			]
	}
		`
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/.well-known/webfinger":
			if r.URL.Query().Get("resource") != "acct:test@"+host {
				t.Errorf("unexpected query value: %s: %s", "resource", r.URL.Query().Get("resource"))
			}

			w.Header().Set("Content-Type", "application/jrd+json")
			io.WriteString(w, renderWebFingerResponse())
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer testServer.Close()

	baseURL = testServer.URL

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	host = u.Host

	client := &webfinger.Client{
		HTTPClient: http.DefaultClient,
		HTTPMode:   true,
	}

	message, err := client.Do(&webfinger.Request{Host: host, Resource: "acct:test@" + host})
	if err != nil {
		t.Error(err)
		return
	}

	if message.Subject != "acct:test@"+host {
		t.FailNow()
	}
}

func Test_Client_Do_XMLResponse(t *testing.T) {

	var baseURL string
	var host string
	renderWebFingerResponse := func() string {
		return `<?xml version='1.0'?>
		<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<Subject>acct:test@` + host + `</Subject>
		<Alias>` + baseURL + `/@test</Alias>
		<Alias>` + baseURL + `/users/test</Alias>
		<Property type="testtype1">teststring1</Property>
		<Property type="testtype2">teststring2</Property>
		<Property type="testtype3" xsi:nil="true" />
		<Property type="testtype4" nillable="true" />
		<Link rel="http://webfinger.net/rel/profile-page" type="text/html" href="` + baseURL + `/@test"/>
		<Link rel="self" type="application/activity+json" href="` + baseURL + `/users/test"/>
		</XRD>`
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/.well-known/webfinger":
			if r.URL.Query().Get("resource") != "acct:test@"+host {
				t.Errorf("unexpected query value: %s: %s", "resource", r.URL.Query().Get("resource"))
			}

			w.Header().Set("Content-Type", "application/xrd+xml")
			io.WriteString(w, renderWebFingerResponse())
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer testServer.Close()

	baseURL = testServer.URL

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	host = u.Host

	client := &webfinger.Client{
		HTTPClient: http.DefaultClient,
		HTTPMode:   true,
	}

	message, err := client.Do(&webfinger.Request{Host: host, Resource: "acct:test@" + host})
	if err != nil {
		t.Error(err)
		return
	}

	if message.Subject != "acct:test@"+host {
		t.FailNow()
	}
}

func Test_Client_Do_CustomUserAgent(t *testing.T) {

	var baseURL string
	var host string
	renderWebFingerResponse := func() string {
		return `
		{
			"subject": "acct:test@` + host + `",
			"aliases": [
					"` + baseURL + `/@test",
					"` + baseURL + `/users/test"
			],
			"links": [
					{
							"rel": "self",
							"type": "application/activity+json",
							"href": "` + baseURL + `/users/test"
					}
			]
	}
		`
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/.well-known/webfinger":
			if r.URL.Query().Get("resource") != "acct:test@"+host {
				t.Errorf("unexpected query value: %s: %s", "resource", r.URL.Query().Get("resource"))
			}

			if r.Header.Get("User-Agent") != "Test User Agent" {
				t.Errorf("unexpected user agent: %s)", r.Header.Get("User-Agent"))
			}

			w.Header().Set("Content-Type", "application/jrd+json")
			io.WriteString(w, renderWebFingerResponse())
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer testServer.Close()

	baseURL = testServer.URL

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	host = u.Host

	client := &webfinger.Client{
		HTTPClient: http.DefaultClient,
		HTTPMode:   true,
		UserAgent:  "Test User Agent",
	}

	message, err := client.Do(&webfinger.Request{Host: host, Resource: "acct:test@" + host})
	if err != nil {
		t.Error(err)
		return
	}

	if message.Subject != "acct:test@"+host {
		t.FailNow()
	}
}

func Test_Client_Do_NotFound(t *testing.T) {
	var host string

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/.well-known/webfinger":
			w.WriteHeader(404)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	host = u.Host

	client := &webfinger.Client{
		HTTPClient: http.DefaultClient,
		HTTPMode:   true,
	}

	message, err := client.Do(&webfinger.Request{Host: host, Resource: "acct:test@" + host})
	if err != nil {
		if !errors.Is(err, webfinger.ErrResourceNotFound) {
			t.Error(err)
		}
		return
	}

	if message != nil {
		t.FailNow()
	}
}

func Test_Client_Do_Gone(t *testing.T) {
	var host string

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/.well-known/webfinger":
			w.WriteHeader(410)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	host = u.Host

	client := &webfinger.Client{
		HTTPClient: http.DefaultClient,
		HTTPMode:   true,
	}

	message, err := client.Do(&webfinger.Request{Host: host, Resource: "acct:test@" + host})
	if err != nil {
		if !errors.Is(err, webfinger.ErrResourceNotFound) {
			t.Error(err)
		}
		return
	}

	if message != nil {
		t.FailNow()
	}
}

func Test_Client_Do_OtherHTTPStatusError(t *testing.T) {
	var host string

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/.well-known/webfinger":
			w.WriteHeader(500)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	host = u.Host

	client := &webfinger.Client{
		HTTPClient: http.DefaultClient,
		HTTPMode:   true,
	}

	message, err := client.Do(&webfinger.Request{Host: host, Resource: "acct:test@" + host})
	if err != nil {
		var webFingerError *webfinger.Error
		var webFingerResponseStatusError *webfinger.WebFingerResponseStatusError
		if !errors.As(err, &webFingerError) {
			t.Error(err)
		} else if !errors.As(webFingerError, &webFingerResponseStatusError) {
			t.Error(err)
		}
		return
	}

	if message != nil {
		t.FailNow()
	}
}

func Test_Client_Do_UnsupportedContentType(t *testing.T) {
	var host string

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/.well-known/webfinger":
			w.Header().Add("Content-Type", "text/plain")
			w.Write([]byte("test"))
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	host = u.Host

	client := &webfinger.Client{
		HTTPClient: http.DefaultClient,
		HTTPMode:   true,
	}

	message, err := client.Do(&webfinger.Request{Host: host, Resource: "acct:test@" + host})
	if err != nil {
		var webFingerError *webfinger.Error
		var unsupportedContentTypeError *webfinger.UnsupportedContentTypeError
		if !errors.As(err, &webFingerError) {
			t.Error(err)
		} else if !errors.As(webFingerError, &unsupportedContentTypeError) {
			t.Error(err)
		}
		return
	}

	if message != nil {
		t.FailNow()
	}
}

func Test_Client_Do_DoHTTPRequestFailed(t *testing.T) {
	var host string

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/.well-known/webfinger":
			w.WriteHeader(410)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	host = u.Host

	client := &webfinger.Client{
		HTTPClient: http.DefaultClient,
		HTTPMode:   true,
	}

	message, err := client.Do(&webfinger.Request{Host: "127.0.0.1.test", Resource: "acct:test@" + host})
	if err != nil {
		var webFingerError *webfinger.Error
		if !errors.As(err, &webFingerError) {
			t.Error(err)
		}
		return
	}

	if message != nil {
		t.FailNow()
	}
}

func Test_Client_Do_InvalidContentType(t *testing.T) {
	var host string

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/.well-known/webfinger":
			w.Header().Add("Content-Type", "multipart/form-data; test=")
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	host = u.Host

	client := &webfinger.Client{
		HTTPClient: http.DefaultClient,
		HTTPMode:   true,
	}

	message, err := client.Do(&webfinger.Request{Host: "127.0.0.1.test", Resource: "acct:test@" + host})
	if err != nil {
		var webFingerError *webfinger.Error
		if !errors.As(err, &webFingerError) {
			t.Error(err)
		}
		return
	}

	if message != nil {
		t.FailNow()
	}
}
