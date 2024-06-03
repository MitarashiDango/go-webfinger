package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/MitarashiDango/ohagi-go-webfinger/nullable"
)

func Test_String_IsZero_ValidString(t *testing.T) {
	s := nullable.String{
		Valid:  true,
		String: "test",
	}

	if s.IsZero() {
		t.Fail()
	}
}

func Test_String_IsZero_InvalidString(t *testing.T) {
	s := nullable.String{
		Valid:  false,
		String: "test",
	}

	if !s.IsZero() {
		t.Fail()
	}
}

func Test_String_SetValue_001(t *testing.T) {
	s := nullable.String{
		Valid:  true,
		String: "",
	}

	s.SetValue("test")

	if s.StringOrZero() != "test" {
		t.Fail()
	}
}

func Test_String_SetValue_002(t *testing.T) {
	s := nullable.String{
		Valid:  false,
		String: "",
	}

	s.SetValue("test")

	if s.StringOrZero() != "test" {
		t.Fail()
	}
}

func Test_String_SetNil_001(t *testing.T) {
	s := nullable.String{
		Valid:  true,
		String: "test",
	}

	s.SetNil()

	if s.StringOrZero() != "" {
		t.Fail()
	}
}

func Test_String_SetNil_002(t *testing.T) {
	s := nullable.String{
		Valid:  false,
		String: "test",
	}

	s.SetNil()

	if s.StringOrZero() != "" {
		t.Fail()
	}
}

func Test_String_StringOrZero_ValidString(t *testing.T) {
	s := nullable.String{
		Valid:  true,
		String: "test",
	}

	if s.StringOrZero() != "test" {
		t.Fail()
	}
}

func Test_String_StringOrZero_InvalidString(t *testing.T) {
	s := nullable.String{
		Valid:  false,
		String: "test",
	}

	if s.StringOrZero() == "test" {
		t.Fail()
	}
}

func Test_String_UnmarshalJSON_001(t *testing.T) {
	jsonString := `{"test":"teststring"}`

	var dest struct {
		Test nullable.String `json:"test"`
	}

	if err := json.Unmarshal([]byte(jsonString), &dest); err != nil {
		t.Fatal(err)
	}

	if dest.Test.StringOrZero() != "teststring" {
		t.Fail()
	}
}

func Test_String_UnmarshalJSON_002(t *testing.T) {
	jsonString := `{"test":null}`

	var dest struct {
		Test nullable.String `json:"test"`
	}

	if err := json.Unmarshal([]byte(jsonString), &dest); err != nil {
		t.Fatal(err)
	}

	if !dest.Test.IsZero() {
		t.Fail()
	}

	if dest.Test.StringOrZero() != "" {
		t.Fail()
	}
}

func Test_String_UnmarshalJSON_003(t *testing.T) {
	jsonString := `{"test":"null"}`

	var dest struct {
		Test nullable.String `json:"test"`
	}

	if err := json.Unmarshal([]byte(jsonString), &dest); err != nil {
		t.Fatal(err)
	}

	if dest.Test.StringOrZero() != "null" {
		t.Fail()
	}
}

func Test_String_UnmarshalJSON_004(t *testing.T) {
	jsonString := `{"test":"nil"}`

	var dest struct {
		Test nullable.String `json:"test"`
	}

	if err := json.Unmarshal([]byte(jsonString), &dest); err != nil {
		t.Fatal(err)
	}

	if dest.Test.StringOrZero() != "nil" {
		t.Fail()
	}
}

func Test_String_MarshalJSON_001(t *testing.T) {
	expected := `{"test":"teststring"}`

	var src struct {
		Test nullable.String `json:"test"`
	}

	src.Test.SetValue("teststring")

	actual, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) != expected {
		t.Fail()
	}
}

func Test_String_MarshalJSON_002(t *testing.T) {
	expected := `{"test":null}`

	var src struct {
		Test nullable.String `json:"test"`
	}

	src.Test.SetNil()

	actual, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) != expected {
		t.Fail()
	}
}

func Test_String_MarshalJSON_003(t *testing.T) {
	expected := `{"test":"null"}`

	var src struct {
		Test nullable.String `json:"test"`
	}

	src.Test.SetValue("null")

	actual, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) != expected {
		t.Fail()
	}
}

func Test_String_MarshalJSON_004(t *testing.T) {
	expected := `{"test":"nil"}`

	var src struct {
		Test nullable.String `json:"test"`
	}

	src.Test.SetValue("nil")

	actual, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) != expected {
		t.Fail()
	}
}
