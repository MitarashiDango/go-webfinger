package nullable

import (
	"encoding/json"
	"fmt"
)

type String struct {
	Valid  bool
	String string
}

func (s String) IsZero() bool {
	return !s.Valid
}

func (s *String) SetValue(v string) {
	s.String, s.Valid = v, true
}

func (s *String) SetNil() {
	s.String, s.Valid = "", false
}

func (s String) Equal(value String) bool {
	return s.Valid == value.Valid && s.String == value.String
}

func (s String) StringOrZero() string {
	if s.Valid {
		return s.String
	}

	return ""
}

func (s *String) UnmarshalJSON(data []byte) error {
	var src interface{}
	if err := json.Unmarshal(data, &src); err != nil {
		return err
	}

	switch v := src.(type) {
	case string:
		s.String, s.Valid = v, true
		return nil
	case nil:
		s.String, s.Valid = "", false
		return nil
	default:
		fmt.Println(v)
		return ErrIncorrectValueType
	}
}

func (s String) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}

	return []byte("null"), nil
}