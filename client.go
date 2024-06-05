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

type Client struct {
	HTTPClient *http.Client
	UserAgent  string
	HTTPMode   bool
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

	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return nil, &Error{
			Err: err,
		}
	}

	var webFingerMessage Message
	switch {
	case mimeType == "application/xrd+xml" || mimeType == "application/xml" || mimeType == "text/xml":
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

	case mimeType == "application/jrd+json" || mimeType == "application/json":
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
				ContentType: mimeType,
			},
		}
	}
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
