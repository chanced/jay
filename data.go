package jay

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Data []byte

func (d Data) Len() int {
	return len(d)
}

func (d Data) MarshalJSON() ([]byte, error) {
	if d == nil {
		return Null, nil
	}

	return d, nil
}

func (d *Data) UnmarshalJSON(data []byte) error {
	if d == nil {
		return errors.New("jay: UnmarshalJSON on nil pointer")
	}
	*d = append((*d)[0:0], data...)
	return nil
}

func (d Data) IsObject() bool {
	return IsObject(d)
}

func (d Data) IsEmptyObject() bool {
	return IsEmptyObject(d) && isEmpty(d)
}

func (d Data) IsEmptyArray() bool {
	return IsEmptyArray(d)
}

// IsArray reports whether the data is a json array. It does not check whether
// the json is malformed.
func (d Data) IsArray() bool {
	return IsArray(d)
}

func (d Data) IsNull() bool {
	return IsNull(d)
}

// IsBool reports true if data appears to be a json boolean value. It is
// possible that it will report false positives of malformed json.
//
// IsBool does not parse strings
func (d Data) IsBool() bool {
	return d.IsTrue() || d.IsFalse()
}

// IsTrue reports true if data appears to be a json boolean value of true. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsTrue does not parse strings
func (d Data) IsTrue() bool {
	return IsTrue(d)
}

// IsFalse reports true if data appears to be a json boolean value of false. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsFalse does not parse strings
func (d Data) IsFalse() bool {
	return IsFalse(d)
}

func (d Data) Equal(data []byte) bool {
	return bytes.Equal(d, data)
}

// ContainsEscapeRune reports whether the string value of d contains "\"
// It returns false if d is not a quoted string.
func (d Data) ContainsEscapeRune() bool {
	for i := 0; i < len(d); i++ {
		if d[i] == '\\' {
			return true
		}
	}
	return false
}

func (d Data) IsNumber() bool {
	return IsNumber(d)
}

func (d Data) IsString() bool {
	return IsString(d)
}

type Object map[string]Data

func (obj Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]Data(obj))
}

func (obj *Object) UnmarshalJSON(data []byte) error {
	var m map[string]Data
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	*obj = m
	return nil
}

func IsEmptyObject(d []byte) bool {
	return IsObject(bytes.TrimSpace(d)) && len(d) == 2
}

// IsValid reports whether s is a valid JSON number literal.
//
// Taken from encoding/json/scanner.go
func IsNumber(data []byte) bool {
	// This function implements the JSON numbers grammar.
	// See https://tools.ietf.org/html/rfc7159#section-6
	// and https://www.json.org/img/number.png

	if len(data) == 0 {
		return false
	}

	// Optional -
	if data[0] == '-' {
		data = data[1:]
		if len(data) == 0 {
			return false
		}
	}

	// Digits
	switch {
	default:
		return false

	case data[0] == '0':
		data = data[1:]

	case '1' <= data[0] && data[0] <= '9':
		data = data[1:]
		for len(data) > 0 && '0' <= data[0] && data[0] <= '9' {
			data = data[1:]
		}
	}

	// . followed by 1 or more digits.
	if len(data) >= 2 && data[0] == '.' && '0' <= data[1] && data[1] <= '9' {
		data = data[2:]
		for len(data) > 0 && '0' <= data[0] && data[0] <= '9' {
			data = data[1:]
		}
	}

	// e or E followed by an optional - or + and
	// 1 or more digits.
	if len(data) >= 2 && (data[0] == 'e' || data[0] == 'E') {
		data = data[1:]
		if data[0] == '+' || data[0] == '-' {
			data = data[1:]
			if len(data) == 0 {
				return false
			}
		}
		for len(data) > 0 && '0' <= data[0] && data[0] <= '9' {
			data = data[1:]
		}
	}

	// Make sure we are at the end.
	return len(data) == 0
}
