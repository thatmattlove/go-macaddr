package macaddr

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
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

// MaskFromPrefixLen creates a MACAddress mask from a prefix bit length. For example, a prefix
// length of 24 would return MAC Address ff:ff:ff:00:00:00.
func MaskFromPrefixLen(l int) *MACAddress {
	if l > _macBitLen {
		return &MACAddress{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	}
	bs := strings.Repeat("1", l) + strings.Repeat("0", _macBitLen-l)
	var ba []byte
	var str string

	for i := len(bs); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(bs[0:i])
		} else {
			str = string(bs[i-8 : i])
		}
		v, _ := strconv.ParseUint(str, 2, 8)
		ba = append([]byte{byte(v)}, ba...)
	}
	return FromByteArray(ba)
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
func (m *MACAddress) Int() int64 {
	if m == nil {
		return 0
	}
	return byteArrayToInt(*m)
}

// ByteString returns a string representation of each MAC Address byte.
func (m *MACAddress) ByteString() string {
	if m == nil {
		return _nilStr
	}
	bsa := []string{}
	for _, b := range *m {
		bsa = append(bsa, fmt.Sprint(b))
	}
	return fmt.Sprintf("{%s}", strings.Join(bsa, ","))
}

func (m *MACAddress) Clone() *MACAddress {
	if m == nil {
		return nil
	}
	return FromBytes((*m)[0], (*m)[1], (*m)[2], (*m)[3], (*m)[4], (*m)[5])
}

// move goes forward or backwards one address from the current. To go backwards, use -1.
func (m *MACAddress) move(p int) *MACAddress {
	if m == nil {
		return nil
	}
	neg := p < 0
	xm := *m.Clone()

	for i := len(xm) - 1; i >= 0; i-- {
		if neg {
			xm[i]--
		} else {
			xm[i]++
		}
		if xm[i] > 0 {
			return FromByteArray(xm)
		}
	}
	return &xm
}

// Next returns the next MACAddress after the current MACAddress.
func (m *MACAddress) Next() *MACAddress {
	return m.move(1)
}

// Previous returns the previous MACAddress before the current MACAddress.
func (m *MACAddress) Previous() *MACAddress {
	return m.move(-1)
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
func (m *MACAddress) Equal(o *MACAddress) bool {
	if m == nil || o == nil {
		return false
	}
	c := bytes.Compare(*m, *o)
	return c == 0
}

// GreaterThan determines if this MACAddress is greater than an input MACAddress.
func (m *MACAddress) GreaterThan(o *MACAddress) bool {
	if m == nil || o == nil {
		return false
	}
	thisI := m.Int()
	thatI := o.Int()
	return thisI > thatI
}

// LessThan determines if this MACAddress is less than an input MACAddress.
func (m *MACAddress) LessThan(o *MACAddress) bool {
	if m == nil || o == nil {
		return false
	}
	thisI := m.Int()
	thatI := o.Int()
	return thisI < thatI
}

// GEqual determines if this MACAddress is greater than or equal to an input MACAddress.
func (m *MACAddress) GEqual(o *MACAddress) bool {
	g := m.GreaterThan(o)
	e := m.Equal(o)
	return g || e
}

// GEqual determines if this MACAddress is less than or equal to an input MACAddress.
func (m *MACAddress) LEqual(o *MACAddress) bool {
	l := m.LessThan(o)
	e := m.Equal(o)
	return l || e
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

// OUI returns the Organizationally Unique Identifier (OUI) of a MACAddress. If a prefix length is
// provided, the MACAddress will be masked with this prefix length. If no prefix length is
// provided, a 24 bit length is assumed.
func (m *MACAddress) OUI(l ...int) string {
	pl := 24
	if len(l) > 0 {
		pl = l[0]
	}
	if m == nil {
		return _nilStr
	}
	mm := m.Mask(MaskFromPrefixLen(pl))

	if pl <= 24 {
		s := mm.String()
		return s[:_hexStrWithColonsLen/2]
	}
	return fmt.Sprintf("%s/%d", mm.String(), pl)
}
