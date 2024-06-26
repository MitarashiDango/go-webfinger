package webfinger

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"mime"
	"net/http"
	"net/url"
)

var DefaultClient = &Client{
	HTTPClient: &http.Client{},
}

type AdditionalMediaTypes struct {
	XML  []string
	JSON []string
}

type Client struct {
	HTTPClient           *http.Client
	UserAgent            string
	HTTPMode             bool
	AdditionalMediaTypes *AdditionalMediaTypes
}

func (client *Client) Do(webFingerRequest *Request) (*Message, error) {
	request, err := client.createHTTPRequest(client.HTTPMode, webFingerRequest)
	if err != nil {
		return nil, &Error{
			Err: err,
		}
	}

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return nil, &Error{
			Err: err,
		}
	}
	defer response.Body.Close()

	if err := client.statusCodeToError(response); err != nil {
		return nil, err
	}

	mediaType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return nil, &Error{
			Err: err,
		}
	}

	additionalMediaTypes := client.AdditionalMediaTypes
	if client.AdditionalMediaTypes == nil {
		additionalMediaTypes = &AdditionalMediaTypes{}
	}

	var webFingerMessage Message
	switch {
	case isXML(mediaType, additionalMediaTypes.XML):
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		if err := xml.Unmarshal(b, &webFingerMessage); err != nil {
			return nil, &Error{
				Err: err,
			}
		}

		return &webFingerMessage, nil

	case isJSON(mediaType, additionalMediaTypes.JSON):
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(b, &webFingerMessage); err != nil {
			return nil, &Error{
				Err: err,
			}
		}

		return &webFingerMessage, nil

	default:
		return nil, &Error{
			Err: &UnsupportedContentTypeError{
				ContentType: mediaType,
			},
		}
	}
}

func isXML(mediaType string, AdditionalMediaTypes []string) bool {
	if mediaType == "application/xrd+xml" || mediaType == "application/xml" || mediaType == "text/xml" {
		return true
	}

	for _, AdditionalMediaType := range AdditionalMediaTypes {
		if mediaType == AdditionalMediaType {
			return true
		}
	}

	return false
}

func isJSON(mediaType string, AdditionalMediaTypes []string) bool {
	if mediaType == "application/jrd+json" || mediaType == "application/json" {
		return true
	}

	for _, AdditionalMediaType := range AdditionalMediaTypes {
		if mediaType == AdditionalMediaType {
			return true
		}
	}

	return false
}

func (client *Client) createHTTPRequest(httpMode bool, webFingerRequest *Request) (*http.Request, error) {
	// requestURL := ?resource=" + url.QueryEscape()
	requestURL, err := url.Parse(getSchema(httpMode) + "//" + webFingerRequest.Host + "/.well-known/webfinger")
	if err != nil {
		return nil, err
	}

	queries := requestURL.Query()
	queries.Set("resource", webFingerRequest.Resource)
	requestURL.RawQuery = queries.Encode()

	request, err := http.NewRequest("GET", requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/jrd+json, application/xrd+xml")

	if client.UserAgent != "" {
		request.Header.Set("User-Agent", client.UserAgent)
	}

	return request, nil
}

func (client *Client) statusCodeToError(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		switch {
		case response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone:
			return ErrResourceNotFound
		default:
			return &Error{
				Err: &WebFingerResponseStatusError{
					URL:        response.Request.URL.String(),
					StatusCode: response.StatusCode,
					Status:     response.Status,
				},
			}
		}
	}

	return nil
}

func getSchema(httpMode bool) string {
	if !httpMode {
		return "https:"
	}

	return "http:"
}
