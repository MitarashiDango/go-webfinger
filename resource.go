package webfinger

import (
	"encoding/xml"

	"github.com/MitarashiDango/ohagi-go-webfinger/nullable"
)

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

func (r *Resource) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var src struct {
		Subject    string   `xml:"Subject"`
		Aliases    []string `xml:"Alias"`
		Properties []struct {
			Type     string `xml:"type,attr"`
			Nil      bool   `xml:"http://www.w3.org/2001/XMLSchema-instance nil,attr"`
			Nullable bool   `xml:"nillable,attr"`
			Value    string `xml:",chardata"`
		} `xml:"Property"`
		Links []*Link `xml:"Link"`
	}

	if err := d.DecodeElement(&src, &start); err != nil {
		return err
	}

	properties := map[string]nullable.String{}
	for _, v := range src.Properties {
		var str nullable.String
		if v.Nil || v.Nullable {
			str.SetNil()
		} else {
			str.SetValue(v.Value)
		}
		properties[v.Type] = str
	}

	r.Subject, r.Aliases, r.Properties, r.Links = src.Subject, src.Aliases, properties, src.Links

	return nil
}

func (r Resource) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return nil
}
