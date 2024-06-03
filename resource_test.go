package webfinger

import (
	"encoding/xml"
	"testing"

	"github.com/MitarashiDango/ohagi-go-webfinger/nullable"
)

func Test_Resource_UnmarshalXML_000(t *testing.T) {
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
		t.Fail()
	}

	if dest.Aliases[0] != "http://localhost/@test" {
		t.Fail()
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
		t.Fail()
	}

	if v, ok := rProperties["testtype2"]; !ok || !v {
		t.Fail()
	}

	if v, ok := rProperties["testtype3"]; !ok || !v {
		t.Fail()
	}

	if v, ok := rProperties["testtype4"]; !ok || !v {
		t.Fail()
	}

	if dest.Links[0].Type != "text/html" {
		t.Fail()
	}

	if dest.Links[0].Rel != "http://webfinger.net/rel/profile-page" {
		t.Fail()
	}

	if dest.Links[0].Type != "text/html" {
		t.Fail()
	}

	if dest.Links[0].Href != "http://localhost/@test" {
		t.Fail()
	}

	if dest.Links[1].Rel != "self" {
		t.Fail()
	}

	if dest.Links[1].Type != "application/activity+json" {
		t.Fail()
	}

	if dest.Links[1].Href != "http://localhost/users/test" {
		t.Fail()
	}
}
