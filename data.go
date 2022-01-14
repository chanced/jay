package jay

import (
	"bytes"
	"encoding/json"
	"errors"
	"unicode"
)

var Null = Data("null")

type Data []byte

func (d Data) Len() int {
	return len(d)
}

func (d Data) MarshalJSON() ([]byte, error) {
	if d == nil || len(d) == 0 {
		return Null, nil
	}

	return d, nil
}

func (d *Data) UnmarshalJSON(data []byte) error {
	if d == nil {
		return errors.New("dynamic.RawMessage: UnmarshalJSON on nil pointer")
	}
	*d = append((*d)[0:0], data...)
	return nil
}

func (d Data) IsObject() bool {
	for _, v := range d {
		if !unicode.IsSpace(rune(v)) {
			return v == '{'
		}
	}
	return false
}

func (d Data) IsEmptyObject() bool {
	return d.IsObject() && len(d) == 2
}

func (d Data) IsEmptyArray() bool {
	if !d.IsArray() {
		return false
	}
	count := 0
	for _, v := range d {
		if !unicode.IsSpace(rune(v)) {
			count += 1
			if count > 2 {
				return false
			}
		}
	}
	return count == 2
}

// IsArray reports whether the data is a json array. It does not check whether
// the json is malformed.
func (d Data) IsArray() bool {
	for _, v := range d {
		if !unicode.IsSpace(rune(v)) {
			return v == '['
		}
	}
	return false
}

func (d Data) IsNull() bool {
	return bytes.Equal(d, Null)
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
	return len(d) == 4 && d[0] == 't'
}

// IsFalse reports true if data appears to be a json boolean value of false. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsFalse does not parse strings
func (d Data) IsFalse() bool {
	return len(d) == 5 && d[0] == 'f'
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

// UnquotedString trims double quotes from the bytes. It does not parse for
// escaped characters
func (d Data) UnquotedString() string {
	if len(d) < 2 {
		return string(d)
	}

	if d[0] == '"' && d[len(d)-1] == '"' {
		return string(d[1 : len(d)-1])
	}
	return string(d)
}

func (d Data) IsNumber() bool {
	if len(d) == 0 {
		return false
	}
	switch d[0] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return true
	default:
		return false
	}
}

func (d Data) IsString() bool {
	if len(d) == 0 {
		return false
	}
	return d[0] == '"'
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
