package webfinger

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/MitarashiDango/ohagi-go-webfinger/nullable"
)

func Test_Resource_UnmarshalXML_001(t *testing.T) {
	xmlString := `<?xml version='1.0'?>
<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<Subject>acct:test@localhost</Subject>
<Alias>http://localhost/@test</Alias>
<Alias>http://localhost/users/test</Alias>
<Property type="testtype1">teststring1</Property>
<Property type="testtype2">teststring2</Property>
<Property type="testtype3" xsi:nil="true" />
<Property type="testtype4" nillable="true" />
<Link rel="http://webfinger.net/rel/profile-page" type="text/html" href="http://localhost/@test"/>
<Link rel="self" type="application/activity+json" href="http://localhost/users/test"/>
</XRD>`

	var dest Resource

	if err := xml.Unmarshal([]byte(xmlString), &dest); err != nil {
		t.Fatal(err)
	}

	if dest.Subject != "acct:test@localhost" {
		t.FailNow()
	}

	if dest.Aliases[0] != "http://localhost/@test" {
		t.FailNow()
	}

	rProperties := map[string]bool{}
	for k, actual := range dest.Properties {
		var expected nullable.String
		switch k {
		case "testtype1":
			expected.SetValue("teststring1")
		case "testtype2":
			expected.SetValue("teststring2")
		case "testtype3":
			expected.SetNil()
		case "testtype4":
			expected.SetNil()
		}

		rProperties[k] = expected.Equal(actual)
	}

	if v, ok := rProperties["testtype1"]; !ok || !v {
		t.FailNow()
	}

	if v, ok := rProperties["testtype2"]; !ok || !v {
		t.FailNow()
	}

	if v, ok := rProperties["testtype3"]; !ok || !v {
		t.FailNow()
	}

	if v, ok := rProperties["testtype4"]; !ok || !v {
		t.FailNow()
	}

	if dest.Links[0].Type != "text/html" {
		t.FailNow()
	}

	if dest.Links[0].Rel != "http://webfinger.net/rel/profile-page" {
		t.FailNow()
	}

	if dest.Links[0].Type != "text/html" {
		t.FailNow()
	}

	if dest.Links[0].Href != "http://localhost/@test" {
		t.FailNow()
	}

	if dest.Links[1].Rel != "self" {
		t.FailNow()
	}

	if dest.Links[1].Type != "application/activity+json" {
		t.FailNow()
	}

	if dest.Links[1].Href != "http://localhost/users/test" {
		t.FailNow()
	}
}

func Test_Resource_UnmarshalJSON_001(t *testing.T) {
	jsonString := `{"subject":"acct:test@localhost","aliases":["http://localhost/@test","http://localhost/users/test"],"properties":{"testtype1":"teststring1","testtype2":"teststring2","testtype3":null,"testtype4":null},"links":[{"rel":"http://webfinger.net/rel/profile-page","type":"text/html","href":"http://localhost/@test"},{"rel":"self","type":"application/activity+json","href":"http://localhost/users/test"}]}`

	var dest Resource

	if err := json.Unmarshal([]byte(jsonString), &dest); err != nil {
		t.Fatal(err)
	}

	if dest.Subject != "acct:test@localhost" {
		t.FailNow()
	}

	if dest.Aliases[0] != "http://localhost/@test" {
		t.FailNow()
	}

	rProperties := map[string]bool{}
	for k, actual := range dest.Properties {
		var expected nullable.String
		switch k {
		case "testtype1":
			expected.SetValue("teststring1")
		case "testtype2":
			expected.SetValue("teststring2")
		case "testtype3":
			expected.SetNil()
		case "testtype4":
			expected.SetNil()
		}

		rProperties[k] = expected.Equal(actual)
	}

	if v, ok := rProperties["testtype1"]; !ok || !v {
		t.FailNow()
	}

	if v, ok := rProperties["testtype2"]; !ok || !v {
		t.FailNow()
	}

	if v, ok := rProperties["testtype3"]; !ok || !v {
		t.FailNow()
	}

	if v, ok := rProperties["testtype4"]; !ok || !v {
		t.FailNow()
	}

	if dest.Links[0].Type != "text/html" {
		t.FailNow()
	}

	if dest.Links[0].Rel != "http://webfinger.net/rel/profile-page" {
		t.FailNow()
	}

	if dest.Links[0].Type != "text/html" {
		t.FailNow()
	}

	if dest.Links[0].Href != "http://localhost/@test" {
		t.FailNow()
	}

	if dest.Links[1].Rel != "self" {
		t.FailNow()
	}

	if dest.Links[1].Type != "application/activity+json" {
		t.FailNow()
	}

	if dest.Links[1].Href != "http://localhost/users/test" {
		t.FailNow()
	}
}

func Test_Resource_MarshalXML_001(t *testing.T) {
	expected := `<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><Subject>acct:test@localhost</Subject><Alias>http://localhost/@test</Alias><Alias>http://localhost/users/test</Alias><Property type="testtype1">teststring1</Property><Property type="testtype2">teststring2</Property><Property type="testtype3" xsi:nil="true"></Property><Link rel="http://webfinger.net/rel/profile-page" type="text/html" href="http://localhost/@test"></Link><Link rel="self" type="application/activity+json" href="http://localhost/users/test"></Link></XRD>`
	resource := &Resource{
		Subject: "acct:test@localhost",
		Aliases: []string{
			"http://localhost/@test",
			"http://localhost/users/test",
		},
		Properties: map[string]nullable.String{
			"testtype1": {
				Valid:  true,
				String: "teststring1",
			},
			"testtype2": {
				Valid:  true,
				String: "teststring2",
			},
			"testtype3": {
				Valid:  false,
				String: "",
			},
		},
		Links: []Link{
			{
				Rel:  "http://webfinger.net/rel/profile-page",
				Type: "text/html",
				Href: "http://localhost/@test",
			},
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: "http://localhost/users/test",
			},
		},
	}

	b, err := xml.Marshal(resource)
	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.FailNow()
	}
}

func Test_Resource_MarshalJSON_001(t *testing.T) {
	expected := `{"subject":"acct:test@localhost","aliases":["http://localhost/@test","http://localhost/users/test"],"properties":{"testtype1":"teststring1","testtype2":"teststring2","testtype3":null},"links":[{"rel":"http://webfinger.net/rel/profile-page","type":"text/html","href":"http://localhost/@test"},{"rel":"self","type":"application/activity+json","href":"http://localhost/users/test"}]}`
	resource := &Resource{
		Subject: "acct:test@localhost",
		Aliases: []string{
			"http://localhost/@test",
			"http://localhost/users/test",
		},
		Properties: map[string]nullable.String{
			"testtype1": {
				Valid:  true,
				String: "teststring1",
			},
			"testtype2": {
				Valid:  true,
				String: "teststring2",
			},
			"testtype3": {
				Valid:  false,
				String: "",
			},
		},
		Links: []Link{
			{
				Rel:  "http://webfinger.net/rel/profile-page",
				Type: "text/html",
				Href: "http://localhost/@test",
			},
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: "http://localhost/users/test",
			},
		},
	}

	b, err := json.Marshal(resource)
	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.FailNow()
	}
}
