package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/thatmattlove/go-macaddr/internal/constant"
)

// Hex ensures all alphanumeric characters in a string are valid hexadecimal characters.
// For example, "abcdef" would return true, but "abcdefg" would return false.
func Hex(i string) (o bool) {
	p := regexp.MustCompile(`[a-zA-Z0-9]+`)
	hp := regexp.MustCompile(`[a-f0-9]+`)
	ss := p.FindAllString(i, -1)
	sj := strings.ToLower(strings.Join(ss, ""))
	hps := strings.Join(hp.FindAllString(sj, -1), "")
	return len(sj) == len(hps)
}

// ParseMacAddrWithPrefixLen operates similarly to ParseMACPrefix, however, it returns the
// validated MAC address object and the prefix length as an integer. If no prefix is provided,
// a /48 prefix length is assumed.
func ParseMacAddrWithPrefixLen(s string) (string, int, error) {
	if !Hex(s) {
		err := fmt.Errorf("'%v' is an invalid MAC address or prefix", s)
		return "", 0, err
	}
	i := strings.IndexByte(s, '/')
	a := s
	if i < 0 {
		i = constant.MacBitLen
	} else {
		aa, ii := s[:i], s[i+1:]
		iii, err := strconv.Atoi(ii)
		if err != nil {
			iii = constant.MacBitLen
		}
		a = aa
		i = iii
	}
	return a, i, nil
}
