package macaddr

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// MACAddress represents a single MAC Address, a slice of bytes. Currently, only 48-bit (EUI-48)
// addresses are supported.
type MACAddress []byte

// ParseMACAddress parses an input string to a valid MACAddress object.
func ParseMACAddress(i string) (o *MACAddress, err error) {
	if !validateHex(i) {
		return nil, fmt.Errorf("'%v' contains non-hexadecimal characters", i)
	}
	hw, err := net.ParseMAC(i)
	if err != nil {
		hw, err = net.ParseMAC(withColons(padMAC(i)))
		if err != nil {
			return nil, err
		}
	}
	o = FromByteArray([]byte{hw[0], hw[1], hw[2], hw[3], hw[4], hw[5]})
	return
}

// MustParseMACAddress operates identically to ParseMACAddress, but panics on error instead of
// returning the error. Most ideal for tests.
func MustParseMACAddress(i string) (o *MACAddress) {
	o, err := ParseMACAddress(i)
	if err != nil {
		panic(err)
	}
	return
}

// Create a MACAddress object directly from bytes.
func FromBytes(one, two, three, four, five, six byte) (m *MACAddress) {
	mac := make(MACAddress, _macByteLen)
	for i, b := range []byte{one, two, three, four, five, six} {
		mac[i] = b
	}
	return &mac
}

// Create a MACAddress object directly from a byte array.
func FromByteArray(b []byte) (m *MACAddress) {
	return FromBytes(b[0], b[1], b[2], b[3], b[4], b[5])
}

// String formats the MAC Address with colons, e.g. 'xx:xx:xx:xx:xx:xx'.
func (m *MACAddress) String() string { return m.Format(_fmtColon) }

// Dots formats the MAC Address with dots, e.g. 'xxxx.xxxx.xxxx'.
func (m *MACAddress) Dots() string { return m.Format(_fmtDot) }

// Dashes formats the MAC Address with dashes, e.g. 'xx-xx-xx-xx-xx-xx'.
func (m *MACAddress) Dashes() string { return m.Format(_fmtDash) }

// NoSeparators formats the MAC Address with no separators, e.g. 'xx-xx-xx-xx-xx-xx'.
func (m *MACAddress) NoSeparators() string { return m.Format(_fmtNone) }

// Int returns an integer representation of a MAC Address.
func (m *MACAddress) Int() int {
	if m == nil {
		return 0
	}
	return byteArrayToInt(*m)
}

// Mask returns the result of masking the MACAddress with the input mask (which is also a
// MACAddress).
func (m *MACAddress) Mask(mask *MACAddress) *MACAddress {
	n := len(*m)
	if n != len(*mask) {
		return nil
	}
	mac := make(MACAddress, n)
	mm := *m
	mp := *mask
	for i := 0; i < n; i++ {
		mac[i] = mm[i] & mp[i]
	}
	return &mac
}

// Equal determines if an input MACAddress is equal to this MACAddress.
func (m *MACAddress) Equal(o *MACAddress) (r bool) {
	c := bytes.Compare(*m, *o)
	return c == 0
}

// Format formats a MACAddress according to a string template. For example, a template of
// xxxx.xxxx.xxxx and a MACAddress of 00:00:5e:00:53:ab would return a value of 0000.5e00.53ab.
func (m *MACAddress) Format(f string) string {
	if m == nil {
		return "<nil>"
	}
	var p []string
	offset := (4 - _macBitLen) & 3
	uc := m.Int() << offset

	fmtStr := createFmtString(reverseString(f))

	for _, ch := range fmtStr {
		if ch == 'x' {
			n := uc & 0xf
			p = append(p, _hexDigits[n])
			uc >>= 4
		} else {
			p = append(p, string(ch))
		}
	}
	return reverseString(strings.Join(p, ""))
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
