package macaddr

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"
)

type MACAddress []byte

func withColons(i string) string {
	p := chunkStr(i, 2)
	r := strings.Join(p, ":")
	return r
}

func padMac(i string) string {
	p := regexp.MustCompile(`[^0-9a-fA-F]+`)
	r := p.ReplaceAllString(i, "")
	return strings.ToLower(padRight(r, "0", _hexStrLen))
}

func formatted(fmtStr string, addr *MACAddress) (f string) {
	if addr == nil {
		return "<nil>"
	}
	var r []string
	offset := (4 - _macBitLen) & 3
	uc := addr.Integer() << offset
	for _, ch := range reverseString(fmtStr) {
		if ch == 'x' {
			n := uc & 0xf
			r = append(r, _hexDigits[n])
			uc >>= 4
		} else {
			r = append(r, string(ch))
		}
	}
	f = reverseString(strings.Join(r, ""))
	return
}

// String formats the MAC Address with colons, e.g. 'xx:xx:xx:xx:xx:xx'.
func (m *MACAddress) String() string { return formatted(_fmtColon, m) }

// Dots formats the MAC Address with dots, e.g. 'xxxx.xxxx.xxxx'.
func (m *MACAddress) Dotted() string { return formatted(_fmtDot, m) }

// Dashes formats the MAC Address with dashes, e.g. 'xx-xx-xx-xx-xx-xx'.
func (m *MACAddress) Dashes() string { return formatted(_fmtDash, m) }

// NoSeparators formats the MAC Address with no separators, e.g. 'xx-xx-xx-xx-xx-xx'.
func (m *MACAddress) NoSeparators() string { return formatted(_fmtNone, m) }

// Integer returns an integer representation of a MAC Address.
func (m *MACAddress) Integer() int {
	if m == nil {
		return 0
	}
	return byteArrayToInt(*m)
}

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

func (m *MACAddress) Equal(o *MACAddress) (r bool) {
	c := bytes.Compare(*m, *o)
	return c == 0
}

func FromBytes(one, two, three, four, five, six byte) (m *MACAddress) {
	mac := make(MACAddress, _macByteLen)
	for i, b := range []byte{one, two, three, four, five, six} {
		mac[i] = b
	}
	return &mac
}

func FromByteArray(b [6]byte) (m *MACAddress) { return FromBytes(b[0], b[1], b[2], b[3], b[4], b[5]) }

func ParseMACAddr(i string) (o *MACAddress, err error) {
	if !validateHex(i) {
		return nil, fmt.Errorf("'%v' contains non-hexadecimal characters", i)
	}
	hw, err := net.ParseMAC(i)
	if err != nil {
		hw, err = net.ParseMAC(withColons(padMac(i)))
		if err != nil {
			return nil, err
		}
	}
	o = FromByteArray([6]byte{hw[0], hw[1], hw[2], hw[3], hw[4], hw[5]})
	return
}

func MustParseMACAddr(i string) (o *MACAddress) {
	o, err := ParseMACAddr(i)
	if err != nil {
		panic(err)
	}
	return
}
