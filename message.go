package webfinger

import (
	"encoding/xml"
	"fmt"
	"slices"

	"github.com/MitarashiDango/ohagi-go-webfinger/nullable"
)

type Properties map[string]nullable.String

type Message struct {
	Subject    string     `json:"subject"`
	Aliases    []string   `json:"aliases,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Links      []Link     `json:"links,omitempty"`
}

type Link struct {
	Rel  string `json:"rel,omitempty" xml:"rel,attr,omitempty"`
	Type string `json:"type,omitempty" xml:"type,attr,omitempty"`
	Href string `json:"href,omitempty" xml:"href,attr,omitempty"`
}

func (r *Message) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var src struct {
		Subject    string   `xml:"Subject"`
		Aliases    []string `xml:"Alias"`
		Properties []struct {
			Type     string `xml:"type,attr"`
			Nil      bool   `xml:"http://www.w3.org/2001/XMLSchema-instance nil,attr"`
			Nullable bool   `xml:"nillable,attr"`
			Value    string `xml:",chardata"`
		} `xml:"Property"`
		Links []Link `xml:"Link"`
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

func (r Message) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type tmpProperties struct {
		Type  string `xml:"type,attr,omitempty"`
		Nil   bool   `xml:"xsi:nil,attr,omitempty"`
		Value string `xml:",chardata"`
	}

	fmt.Println(start)
	start.Name.Local = "XRD"
	start.Name.Space = "http://docs.oasis-open.org/ns/xri/xrd-1.0"
	start.Attr = append(start.Attr, xml.Attr{
		Name: xml.Name{
			Local: "xmlns:xsi",
		},
		Value: "http://www.w3.org/2001/XMLSchema-instance",
	})

	var src struct {
		Subject    string          `xml:"Subject"`
		Aliases    []string        `xml:"Alias,omitempty"`
		Properties []tmpProperties `xml:"Property,omitempty"`
		Links      []Link          `xml:"Link,omitempty"`
	}

	src.Subject = r.Subject
	src.Aliases = r.Aliases
	src.Links = r.Links

	mapKeys := make([]string, 0, len(r.Properties))
	for k := range r.Properties {
		mapKeys = append(mapKeys, k)
	}
	slices.Sort(mapKeys)

	src.Properties = make([]tmpProperties, 0, len(mapKeys))
	for _, k := range mapKeys {
		v := r.Properties[k]

		if v.IsZero() {
			src.Properties = append(src.Properties, tmpProperties{
				Type: k,
				Nil:  true,
			})
		} else {
			src.Properties = append(src.Properties, tmpProperties{
				Type:  k,
				Nil:   false,
				Value: v.StringOrZero(),
			})
		}
	}

	return e.EncodeElement(src, start)
}
