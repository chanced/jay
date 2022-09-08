package jay

import (
	"bytes"
	"encoding/json"
	"errors"
)

type JSON []byte

func (d JSON) Len() int {
	return len(d)
}

func (d JSON) MarshalJSON() ([]byte, error) {
	if d == nil {
		return Null, nil
	}

	return d, nil
}

func (d *JSON) UnmarshalJSON(data []byte) error {
	if d == nil {
		return errors.New("jay: UnmarshalJSON on nil pointer")
	}
	*d = append((*d)[0:0], data...)
	return nil
}

func (d JSON) IsObject() bool {
	return IsObject(d)
}

func (d JSON) IsEmptyObject() bool {
	return IsEmptyObject(d) && isEmpty(d)
}

func (d JSON) IsEmptyArray() bool {
	return IsEmptyArray(d)
}

// IsArray reports whether the data is a json array. It does not check whether
// the json is malformed.
func (d JSON) IsArray() bool {
	return IsArray(d)
}

func (d JSON) IsNull() bool {
	return IsNull(d)
}

// IsBool reports true if data appears to be a json boolean value. It is
// possible that it will report false positives of malformed json.
//
// IsBool does not parse strings
func (d JSON) IsBool() bool {
	return d.IsTrue() || d.IsFalse()
}

// IsTrue reports true if data appears to be a json boolean value of true. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsTrue does not parse strings
func (d JSON) IsTrue() bool {
	return IsTrue(d)
}

// IsFalse reports true if data appears to be a json boolean value of false. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsFalse does not parse strings
func (d JSON) IsFalse() bool {
	return IsFalse(d)
}

func (d JSON) Equal(data []byte) bool {
	return bytes.Equal(d, data)
}

// ContainsEscapeRune reports whether the string value of d contains "\"
// It returns false if d is not a quoted string.
func (d JSON) ContainsEscapeRune() bool {
	for i := 0; i < len(d); i++ {
		if d[i] == '\\' {
			return true
		}
	}
	return false
}

func (d JSON) IsNumber() bool {
	return IsNumber(d)
}

func (d JSON) IsString() bool {
	return IsString(d)
}

type Object map[string]JSON

func (obj Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]JSON(obj))
}

func (obj *Object) UnmarshalJSON(data []byte) error {
	var m map[string]JSON
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
