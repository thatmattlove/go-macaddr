package validate

import (
	"regexp"
	"strings"
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
