package jay

import (
	"bytes"
	"unicode"
)

var (
	Null     = Data("null")
	True     = Data("true")
	False    = Data("false")
	trueVal  = []byte("true")  // unnecessary but what if you accidentally override the above?
	falseVal = []byte("false") // ^^
	nullVal  = []byte("null")
)

func IsNull(d []byte) bool {
	return bytes.Equal(d, nullVal)
}

func IsBool(d []byte) bool {
	return IsTrue(d) || IsFalse(d)
}

func IsString(d []byte) bool {
	return startsAndEndsWith(bytes.TrimSpace(d), '"', '"')
}

func IsObject(d []byte) bool {
	return startsAndEndsWith(bytes.TrimSpace(d), '{', '}')
}

func IsArray(d []byte) bool {
	return startsAndEndsWith(bytes.TrimSpace(d), '[', ']')
}

func IsEmptyArray(d []byte) bool {
	return IsArray(d) && isEmpty(d)
}

// IsTrue reports true if data appears to be a json boolean value of true. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsTrue does not parse strings
func IsTrue(d []byte) bool {
	return bytes.Equal(d, True)
}

// IsFalse reports true if data appears to be a json boolean value of false. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsFalse does not parse strings
func IsFalse(d []byte) bool {
	return bytes.Equal(d, falseVal)
}

func isEmpty(d []byte) bool {
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

func startsAndEndsWith(d []byte, start, end byte) bool {
	if len(d) < 2 {
		return false
	}

	return d[0] == start && d[len(d)-1] == end
}
