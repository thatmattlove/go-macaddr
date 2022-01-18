package macaddr

import (
	"regexp"
	"strings"
)

// byteArrayToInt converts a byte array to an integer.
func byteArrayToInt(arr []byte) int {
	var res int
	for _, v := range arr {
		res <<= 8
		res |= int(v)
	}
	return res
}

// chunkStr chunks a string into chunks of n size. For example, "0123456789ab" with a size of 2
// would become [01 23 45 67 89 ab].
func chunkStr(str string, size int) []string {
	var res []string
	slice := []rune(str)

	if len(slice) == 0 || size <= 0 {
		return res
	}

	length := len(slice)
	if size == 1 || size >= length {
		for _, v := range slice {
			var tmp rune
			tmp += v
			// tmp = append(tmp, v)
			res = append(res, string(tmp))
		}
		return res
	}

	// divide slice equally
	divideNum := length/size + 1
	for i := 0; i < divideNum; i++ {
		if i == divideNum-1 {
			if len(slice[i*size:]) > 0 {
				res = append(res, string(slice[i*size:]))
			}
		} else {
			res = append(res, string(slice[i*size:(i+1)*size]))
		}
	}

	return res
}

// padRight pads a string with a another string up to n total length. For example, given
// arguments:
//
// str: "0123", pad: "0", count: 12
//
// The return value would be "012300000000"
func padRight(str string, pad string, count int) string {
	n := count - len(str)
	if n > 0 {
		return str + strings.Repeat(pad, n)
	}
	return str
}

// DecToInt converts decimal string to integer.
// Returns number, characters consumed, success.
// See Golang internal: https://github.com/golang/go/blob/a59e33224e42d60a97fa720a45e1b74eb6aaa3d0/src/net/parse.go#L122
func decToInt(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= _big {
			return _big, i, false
		}
	}
	if i == 0 {
		return 0, 0, false
	}
	return n, i, true
}

// reverseString reverses a string. For example, "abc" becomes "cba".
func reverseString(i string) (o string) {
	rr := []rune(i)
	var r []rune
	for i := len(rr) - 1; i >= 0; i-- {
		r = append(r, rr[i])
	}
	o = string(r)
	return
}

// reverseString reverses a string. For example, "abc" becomes "cba".
func reverseBytes(i []byte) []byte {
	if len(i) == 0 {
		return i
	}
	o := []byte{}
	for idx := len(i) - 1; idx >= 0; idx-- {
		o = append(o, i[idx])
	}
	return o
}

// validateHex ensures all alphanumeric characters in a string are valid hexadecimal characters.
// For example, "abcdef" would return true, but "abcdefg" would return false.
func validateHex(i string) (o bool) {
	p := regexp.MustCompile(`[a-zA-Z0-9]+`)
	hp := regexp.MustCompile(`[a-f0-9]+`)
	ss := p.FindAllString(i, -1)
	sj := strings.ToLower(strings.Join(ss, ""))
	hps := strings.Join(hp.FindAllString(sj, -1), "")
	return len(sj) == len(hps)
}

// createFmtString parses an input string to replace all alphanumeric characters with 'x', so that
// any valid MAC Address string can be used as a template.
func createFmtString(s string) (f string) {
	p := regexp.MustCompile(`[a-zA-Z0-9]`)
	f = strings.Map(func(r rune) rune {
		if p.MatchString(string(r)) {
			return 'x'
		}
		return r
	}, s)
	return
}

// padMAC right-pads an input string with zeros to guarantee the string length is 12. For example,
// 012345 becomes 012345000000.
func padMAC(i string) string {
	p := regexp.MustCompile(`[^0-9a-fA-F]+`)
	r := p.ReplaceAllString(i, "")
	return strings.ToLower(padRight(r, "0", _hexStrLen))
}

// withColons chunks an input string into n parts of 2 characters, and joins them with colons. For
// example, 0123456789ab becomes 01:23:45:67:89:ab.
func withColons(i string) string {
	p := chunkStr(i, 2)
	return strings.Join(p, ":")
}
