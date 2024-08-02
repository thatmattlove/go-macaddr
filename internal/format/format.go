package format

import (
	"regexp"
	"strings"

	"github.com/thatmattlove/go-macaddr/internal/constant"
	"github.com/thatmattlove/go-macaddr/internal/convert"
)

// PadRight pads a string with a another string up to n total length. For example, given
// arguments:
//
// str: "0123", pad: "0", count: 12
//
// The return value would be "012300000000"
func PadRight(str string, pad string, count int) string {
	n := count - len(str)
	if n > 0 {
		return str + strings.Repeat(pad, n)
	}
	return str
}

// ReverseString reverses a string. For example, "abc" becomes "cba".
func ReverseString(i string) (o string) {
	rr := []rune(i)
	var r []rune
	for i := len(rr) - 1; i >= 0; i-- {
		r = append(r, rr[i])
	}
	o = string(r)
	return
}

// ReverseBytes reverses a byte slice. For example, "abc" becomes "cba".
func ReverseBytes(i []byte) []byte {
	if len(i) == 0 {
		return i
	}
	o := []byte{}
	for idx := len(i) - 1; idx >= 0; idx-- {
		o = append(o, i[idx])
	}
	return o
}

// CreateFmtString parses an input string to replace all alphanumeric characters with 'x', so that
// any valid MAC Address string can be used as a template.
func CreateFmtString(s string) (f string) {
	p := regexp.MustCompile(`[a-zA-Z0-9]`)
	f = strings.Map(func(r rune) rune {
		if p.MatchString(string(r)) {
			return 'x'
		}
		return r
	}, s)
	return
}

// PadMAC right-pads an input string with zeros to guarantee the string length is 12. For example,
// 012345 becomes 012345000000.
func PadMAC(i string) string {
	p := regexp.MustCompile(`[^0-9a-fA-F]+`)
	r := p.ReplaceAllString(i, "")
	return strings.ToLower(PadRight(r, "0", constant.HexStrLen))
}

// WithColons chunks an input string into n parts of 2 characters, and joins them with colons. For
// example, 0123456789ab becomes 01:23:45:67:89:ab.
func WithColons(i string) string {
	p := convert.ChunkStr(i, 2)
	return strings.Join(p, ":")
}
