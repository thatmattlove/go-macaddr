package macaddr

import (
	"fmt"
	"strconv"
	"strings"
)

// MACPrefix represents a MAC Address prefix or range, e.g. 00:00:5e:00:00:00/24 (with a range of
// 00:00:5e:00:00:00-00:00:5e:ff:ff:ff).
type MACPrefix struct {
	// MAC is the base MAC address of a prefix. For example, a MACPrefix of 00:00:5e:00:00:00/24
	// would have a base MAC of 00:00:5e:00:00:00.
	MAC *MACAddress
	// Mask is the mask representation of the MACPrefix. For example, a MACPrefix of
	// 00:00:5e:00:00:00/24 would have a Mask of ff:ff:ff:00:00:00.
	Mask *MACAddress
}

// ParseMACPrefix attempts to parse an input string to a valid MACPrefix object.
//
// Return Values
//
// ParseMACPrefix returns the original input MAC Address as a valid MACAddress object, the
// parsed MACPrefix object, and an error if parsing failed.
func ParseMACPrefix(s string) (mac *MACAddress, mpo *MACPrefix, err error) {
	mac, l, err := parseMacAddrWithPrefixLen(s)
	if err != nil {
		return nil, nil, err
	}
	ls := fmt.Sprint(l)

	n, i, ok := decToInt(ls)
	if mac == nil || !ok || i != len(ls) || n < 0 || n > _macBitLen {
		return nil, nil, fmt.Errorf("'%v' is an invalid MAC prefix", s)
	}
	m := cidrMask(n, _macBitLen)
	var mp *MACPrefix = new(MACPrefix)
	mp.MAC = mac.Mask(m)
	mp.Mask = m
	mpo = mp
	return
}

// MustParseMACPrefix operates identically to ParseMACPrefix, but panics upon parsing error
// instead of returning the error. Most ideal for tests or pre-validated string input.
func MustParseMACPrefix(s string) (mac *MACAddress, mp *MACPrefix) {
	mac, mp, err := ParseMACPrefix(s)
	if err != nil {
		panic(err)
	}
	return
}

// String returns a colon-separated string representation of the MACPrefix object.
func (p *MACPrefix) String() string {
	if p == nil {
		return _nilStr
	}

	if p.MAC == nil || p.Mask == nil {
		return _nilStr
	}
	l := prefixLength(*p.Mask)
	if l == -1 {
		return p.MAC.String() + "/" + p.Mask.String()
	}

	return p.MAC.String() + "/" + strconv.Itoa(l)
}

// Match attempts to match the MACPrefix to an input string.
func (p *MACPrefix) Match(i string) (m *MACPrefix, e error) {
	e = fmt.Errorf("'%v' is not contained within MACPrefix %s", i, p.String())
	addr, l, err := parseMacAddrWithPrefixLen(i)
	if err != nil {
		return nil, err
	}
	pl := p.PrefixLen()

	if l < pl {
		return nil, e
	}

	if p.Contains(addr) {
		return p, nil
	}
	return nil, e
}

// Contains determines if an input MACAddress is contained within this MACPrefix.
func (p *MACPrefix) Contains(mac *MACAddress) bool {
	if p == nil {
		panic(fmt.Errorf("cannot check if MAC '%s' is contained within nil MACPrefix", mac.String()))
	}
	mask := *p.Mask
	if x := mac; x != nil {
		mac = x
	}
	bm := *p.MAC
	im := *mac
	l := len(im)

	if l != len(bm) {
		return false
	}

	for i := 0; i < l; i++ {
		if bm[i]&mask[i] != im[i]&mask[i] {
			return false
		}
	}
	return true
}

// PrefixLen returns the prefix length of the MACPrefix as an integer.
func (p *MACPrefix) PrefixLen() int {
	if p == nil {
		return 0
	}
	return prefixLength(*p.Mask)
}

// OUI returns the Organizationally Unique Identifier (OUI) of a MACPrefix.
func (p *MACPrefix) OUI() string {
	if p == nil {
		return _nilStr
	}
	if p.PrefixLen() <= 24 {
		s := p.String()
		return s[:_hexStrWithColonsLen/2]
	}
	return p.String()
}

// cidrMask returns an MAC Address consisting of 'ones' 1 bits
// followed by 0s up to a total length of 'bits' bits.
// Adapted from: https://github.com/golang/go/blob/2639f2f79bda2c3a4e9ef7381ca7de14935e2a4a/src/net/ip.go#L77
func cidrMask(ones, bits int) *MACAddress {
	if bits != _macBitLen {
		return nil
	}
	if ones < 0 || ones > bits {
		return nil
	}
	m := make(MACAddress, _macByteLen)
	n := uint(ones)
	for i := 0; i < _macByteLen; i++ {
		if n >= 8 {
			m[i] = 0xff
			n -= 8
			continue
		}
		m[i] = ^byte(0xff >> n)
		n = 0
	}
	return &m
}

// prefixLength is adapted from:
// https://github.com/golang/go/blob/2639f2f79bda2c3a4e9ef7381ca7de14935e2a4a/src/net/ip.go#L451
func prefixLength(mac MACAddress) int {
	var n int
	for i, v := range mac {
		if v == 0xff {
			n += 8
			continue
		}
		// found non-ff byte
		// count 1 bits
		for v&0x80 != 0 {
			n++
			v <<= 1
		}
		// rest must be 0 bits
		if v != 0 {
			return -1
		}
		for i++; i < len(mac); i++ {
			if mac[i] != 0 {
				return -1
			}
		}
		break
	}
	return n
}

// parseMacAddrWithPrefixLen operates similarly to ParseMACPrefix, however, it returns the
// validated MAC address object and the prefix length as an integer. If no prefix is provided,
// a /24 prefix length is assumed.
func parseMacAddrWithPrefixLen(s string) (m *MACAddress, l int, err error) {
	if !validateHex(s) {
		err = fmt.Errorf("'%v' is an invalid MAC address or prefix", s)
		return nil, 0, err
	}
	i := strings.IndexByte(s, '/')
	a := s
	if i < 0 {
		i = 24
	} else {
		aa, ii := s[:i], s[i+1:]
		iii, err := strconv.Atoi(ii)

		if err != nil {
			iii = 24
		}
		a = aa
		i = iii
	}
	mac, err := ParseMACAddress(a)
	if err != nil {
		return nil, 0, err
	}

	return mac, i, nil
}
