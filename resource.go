package webfinger

import "github.com/MitarashiDango/ohagi-go-webfinger/nullable"

type Resource struct {
	Subject    string   `json:"subject" xml:"Subject"`
	Aliases    []string `json:"aliases" xml:"Alias"`
	Properties map[string]nullable.String
	Links      []*Link `json:"links" xml:"Link"`
}

type Link struct {
	Rel  string `json:"rel" xml:"rel,attr"`
	Type string `json:"type" xml:"type,attr"`
	Href string `json:"href" xml:"href,attr"`
}
