package webfinger

import "net/http"

type Response struct {
	RawHTTPResponse   *http.Response
	WebFingerResource *Resource
}
