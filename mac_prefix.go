package macaddr

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// MACPrefixIterator tracks iteration state while iterating through a MACPrefix.
type MACPrefixIterator struct {
	err     error
	runs    int
	prefix  *MACPrefix
	last    *MACAddress
	current *MACAddress
}

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
	m := MaskFromPrefixLen(n)
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
		return fmt.Sprintf("%s/%s", p.MAC.String(), p.Mask.String())
	}
	return fmt.Sprintf("%s/%d", p.MAC.String(), l)
}

// Match attempts to match the MACPrefix to an input string.
func (p *MACPrefix) Match(i string) (m *MACPrefix, e error) {
	e = fmt.Errorf("'%v' is not contained within MACPrefix %s", i, p.String())
	addr, l, err := parseMacAddrWithPrefixLen(i)
	if err != nil {
		return nil, err
	}

	if l < p.PrefixLen() {
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

// First returns the first MAC Address in a MACPrefix.
func (p *MACPrefix) First() (mac *MACAddress) {
	if p == nil {
		return nil
	}
	return p.MAC.Clone()
}

// Last returns the last MAC Address in a MACPrefix.
func (p *MACPrefix) Last() (mac *MACAddress) {
	if p == nil {
		return nil
	}
	last := make([]byte, _macByteLen)
	w := *p.WildcardMask()
	for i, b := range *p.MAC {
		last[i] = b + w[i]
	}
	mac = FromBytes(last[0], last[1], last[2], last[3], last[4], last[5])
	return
}

// Count returns the number of MAC Addresses in a MACPrefix.
func (p *MACPrefix) Count() int {
	if p == nil {
		return 0
	}
	exp := _macBitLen - p.PrefixLen()

	if exp == 0 {
		return 1
	}
	c := int(math.Pow(2, float64(exp)))
	return c
}

// WildcardMask returns a MACAddress object of the wildcard mask of the MACPrefix.
func (p *MACPrefix) WildcardMask() (mask *MACAddress) {
	if p == nil {
		return nil
	}
	wc := make([]byte, len(*p.Mask))
	for pos, b := range *p.Mask {
		wc[pos] = 0xff - b
	}
	mask = FromByteArray(wc)
	return
}

// Next iterates through the MACPrefix range.
func (i *MACPrefixIterator) Next() bool {
	if i == nil || i.prefix == nil || i.err != nil {
		return false
	}

	if i.runs > 0 {
		i.current = i.current.Next()
	}
	i.runs++

	return i.current.LEqual(i.last)
}

// Value returns the current iteration value.
func (i *MACPrefixIterator) Value() *MACAddress {
	if i == nil || i.err != nil || i.prefix == nil {
		panic(fmt.Errorf("cannot call Iter() after iterator has finished"))
	}
	return i.current
}

// Iter creates an iterator for the MACPrefix.
func (p *MACPrefix) Iter() *MACPrefixIterator {
	if p == nil {
		return nil
	}
	var err error
	return &MACPrefixIterator{
		prefix:  p,
		last:    p.Last(),
		current: p.First(),
		err:     err,
		runs:    0,
	}
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
// a /48 prefix length is assumed.
func parseMacAddrWithPrefixLen(s string) (m *MACAddress, l int, err error) {
	if !validateHex(s) {
		err = fmt.Errorf("'%v' is an invalid MAC address or prefix", s)
		return nil, 0, err
	}
	i := strings.IndexByte(s, '/')
	a := s
	if i < 0 {
		i = _macBitLen
	} else {
		aa, ii := s[:i], s[i+1:]
		iii, err := strconv.Atoi(ii)

		if err != nil {
			iii = _macBitLen
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
