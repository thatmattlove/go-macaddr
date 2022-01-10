package macaddr

import (
	"fmt"
	"strconv"
	"strings"
)

type MACPrefix struct {
	MAC  MACAddress // base MAC
	Mask MACAddress
}

func (p *MACPrefix) String() string {
	if p == nil {
		return "<nil>"
	}
	bm, m := baseMACAndMask(p)
	if bm == nil || m == nil {
		return "<nil>"
	}
	l := prefixLength(m)
	if l == -1 {
		return bm.String() + "/" + m.String()
	}

	return bm.String() + "/" + strconv.Itoa(l)
}

func (p *MACPrefix) Contains(mac MACAddress) bool {
	bm, mask := baseMACAndMask(p)
	if x := mac; x != nil {
		mac = x
	}
	l := len(mac)
	if l != len(bm) {
		return false
	}
	for i := 0; i < l; i++ {
		if bm[i]&mask[i] != mac[i]&mask[i] {
			return false
		}
	}
	return true
}

func (p *MACPrefix) PrefixLen() int {
	return prefixLength(p.Mask)
}

func ParseMACPrefix(s string) (mac MACAddress, mp *MACPrefix, err error) {
	err = nil
	i := strings.IndexByte(s, '/')
	if i < 0 {
		return nil, nil, fmt.Errorf("'%v' is an invalid MAC prefix", s)
	}
	addr, mask := s[:i], s[i+1:]
	mac, err = ParseMACAddr(addr)
	if err != nil {
		return nil, nil, err
	}

	n, i, ok := decToInt(mask)
	if mac == nil || !ok || i != len(mask) || n < 0 || n > MAC_BIT_LEN {
		return nil, nil, fmt.Errorf("'%v' is an invalid MAC prefix", s)
	}
	m := cidrMask(n, MAC_BIT_LEN)
	mp = &MACPrefix{MAC: mac.Mask(m), Mask: m}
	return mac, mp, err
}

func MustParseMACPrefix(s string) (mac MACAddress, mp *MACPrefix) {
	mac, mp, err := ParseMACPrefix(s)
	if err != nil {
		panic(err)
	}
	return
}

// cidrMask returns an MAC Address consisting of 'ones' 1 bits
// followed by 0s up to a total length of 'bits' bits.
// Adapted from: https://github.com/golang/go/blob/2639f2f79bda2c3a4e9ef7381ca7de14935e2a4a/src/net/ip.go#L77
func cidrMask(ones, bits int) (m MACAddress) {
	if bits != MAC_BIT_LEN {
		return nil
	}
	if ones < 0 || ones > bits {
		return nil
	}
	m = make(MACAddress, MAC_BYTE_LEN)
	n := uint(ones)
	for i := 0; i < MAC_BYTE_LEN; i++ {
		if n >= 8 {
			m[i] = 0xff
			n -= 8
			continue
		}
		m[i] = ^byte(0xff >> n)
		n = 0
	}
	return m
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
