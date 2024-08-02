package macaddr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatmattlove/go-macaddr/internal/constant"
)

func Test_MustParseMACAddress(t *testing.T) {
	t.Run("MustParseMACAddress should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseMACAddress("this should panic")
		})
	})
	t.Run("MustParseMACAddress should not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			// This should not panic.
			MustParseMACAddress("01:23:45:67:89:ab")
		})
	})
}

func Test_ParseMACAddress(t *testing.T) {
	tests := [][]string{
		{"01:23:45:67:89:ab", "01:23:45:67:89:ab"},
		{"01-23-45-67-89-ab", "01:23:45:67:89:ab"},
		{"0123.4567.89.ab", "01:23:45:67:89:ab"},
		{"0123456789ab", "01:23:45:67:89:ab"},
		{"01:23:45", "01:23:45:00:00:00"},
		{"01-23-45-67", "01:23:45:67:00:00"},
	}
	for i, p := range tests {
		t.Run(fmt.Sprintf("ParseMACAddress '%d'", i+1), func(t *testing.T) {
			e := p[1]
			r, err := ParseMACAddress(p[0])
			assert.Nil(t, err)
			assert.Equal(t, r.String(), e)
		})
	}
	t.Run("ParseMACAddress returns error", func(t *testing.T) {
		_, err := ParseMACAddress("0123.4567.89az")
		assert.NotNil(t, err)
	})
}

func Test_MACAddress(t *testing.T) {
	s := "01:23:45:67:89:ab"
	m, err := ParseMACAddress(s)
	assert.Nil(t, err)
	t.Run("Int() returns int", func(t *testing.T) {
		i := m.Int()
		var e int64 = 1250999896491
		assert.Equal(t, e, i)
	})
	t.Run("Int() returns 0", func(t *testing.T) {
		var m *MACAddress
		i := m.Int()
		var e int64 = 0
		assert.Equal(t, e, i)
	})
	t.Run("Mask() works properly", func(t *testing.T) {
		macOut := "01:23:45:00:00:00"
		maskIn := MustParseMACAddress("ff:ff:ff:00:00:00")
		maskOut := m.Mask(maskIn)
		assert.Equal(t, macOut, maskOut.String())
	})
	t.Run("Mask() returns nil", func(t *testing.T) {
		maskIn := MACAddress{}
		maskOut := m.Mask(&maskIn)
		assert.Nil(t, maskOut)
	})
	t.Run("MACAddress.Equal()", func(t *testing.T) {
		et := MustParseMACAddress(s)
		ef := MustParseMACAddress("12:34:56:78:9a:bc")
		assert.True(t, m.Equal(et))
		assert.False(t, m.Equal(ef))
	})
	t.Run("MACAddress.String()", func(t *testing.T) {
		assert.Equal(t, m.String(), s)
	})
	t.Run("MACAddress.String() returns nil", func(t *testing.T) {
		var m *MACAddress
		assert.Equal(t, constant.NilStr, m.String())
	})
	t.Run("MACAddress.Dots()", func(t *testing.T) {
		assert.Equal(t, m.Dots(), "0123.4567.89ab")
	})
	t.Run("MACAddress.Dashes()", func(t *testing.T) {
		assert.Equal(t, m.Dashes(), "01-23-45-67-89-ab")
	})
	t.Run("MACAddress.NoSeparators()", func(t *testing.T) {
		assert.Equal(t, m.NoSeparators(), "0123456789ab")
	})
	t.Run("MACAddress.OUI() 1", func(t *testing.T) {
		oui := m.OUI()
		assert.Equal(t, "01:23:45", oui)
	})
	t.Run("MACAddress.OUI() 2", func(t *testing.T) {
		oui := m.OUI(24)
		assert.Equal(t, "01:23:45", oui)
	})
	t.Run("MACAddress.OUI() 3", func(t *testing.T) {
		m := MustParseMACAddress("00:55:DA:80:01:23")
		oui := m.OUI(28)
		assert.Equal(t, "00:55:da:80:00:00/28", oui)
	})
	t.Run("MACAddress.OUI() nil", func(t *testing.T) {
		var m *MACAddress
		oui := m.OUI()
		assert.Equal(t, constant.NilStr, oui)
	})
	t.Run("MACAddress.ByteString() nil", func(t *testing.T) {
		var m *MACAddress
		assert.Equal(t, constant.NilStr, m.ByteString())
	})
	t.Run("MACAddress.ByteString() 1", func(t *testing.T) {
		e := "{1,35,69,103,137,171}"
		assert.Equal(t, e, m.ByteString())
	})
	t.Run("MACAddress.move() nil", func(t *testing.T) {
		var m *MACAddress
		assert.Nil(t, m.move(1))
	})
	t.Run("MACAddress.Next()", func(t *testing.T) {
		e := MustParseMACAddress("01:23:45:67:89:ac")
		assert.Equal(t, e, m.Next())
	})
	t.Run("MACAddress.Previous()", func(t *testing.T) {
		e := MustParseMACAddress("01:23:45:67:89:aa")
		assert.Equal(t, e, m.Previous())
	})
	t.Run("MACAddress.GreaterThan() nil", func(t *testing.T) {
		var m *MACAddress
		mm := FromBytes(0x01, 0x23, 0x45, 0x67, 0x89, 0xab)
		assert.False(t, m.GreaterThan(&MACAddress{}))
		assert.False(t, mm.GreaterThan(nil))
	})
	t.Run("MACAddress.GreaterThan() 1", func(t *testing.T) {
		e := &MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaa}
		assert.True(t, m.GreaterThan(e))
	})
	t.Run("MACAddress.GreaterThan() 2", func(t *testing.T) {
		e := &MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaf}
		assert.False(t, m.GreaterThan(e))
	})
	// LessThan()
	t.Run("MACAddress.LessThan() nil", func(t *testing.T) {
		var m *MACAddress
		mm := FromBytes(0x01, 0x23, 0x45, 0x67, 0x89, 0xab)
		assert.False(t, m.LessThan(&MACAddress{}))
		assert.False(t, mm.LessThan(nil))
	})
	t.Run("MACAddress.LessThan() 1", func(t *testing.T) {
		e := &MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaf}
		assert.True(t, m.LessThan(e))
	})
	t.Run("MACAddress.LessThan() 2", func(t *testing.T) {
		e := &MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaa}
		assert.False(t, m.LessThan(e))
	})
	// GEqual()
	t.Run("MACAddress.GEqual() nil", func(t *testing.T) {
		var m *MACAddress
		mm := FromBytes(0x01, 0x23, 0x45, 0x67, 0x89, 0xab)
		assert.False(t, m.GEqual(&MACAddress{}))
		assert.False(t, mm.GEqual(nil))
	})
	t.Run("MACAddress.GEqual() 1", func(t *testing.T) {
		e := &MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaa}
		assert.True(t, m.GEqual(e))
	})
	t.Run("MACAddress.GEqual() 2", func(t *testing.T) {
		e := &MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaf}
		assert.False(t, m.GEqual(e))
	})
	t.Run("MACAddress.GEqual() 3", func(t *testing.T) {
		e := &MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xab}
		assert.True(t, m.GEqual(e))
	})
	// LEqual()
	t.Run("MACAddress.LEqual() nil", func(t *testing.T) {
		var m *MACAddress
		mm := FromBytes(0x01, 0x23, 0x45, 0x67, 0x89, 0xab)
		assert.False(t, m.LEqual(&MACAddress{}))
		assert.False(t, mm.LEqual(nil))
	})
	t.Run("MACAddress.LEqual() 1", func(t *testing.T) {
		e := &MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaf}
		assert.True(t, m.LEqual(e))
	})
	t.Run("MACAddress.LEqual() 2", func(t *testing.T) {
		e := &MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaa}
		assert.False(t, m.LEqual(e))
	})
	t.Run("MACAddress.LEqual() 3", func(t *testing.T) {
		e := &MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xab}
		assert.True(t, m.LEqual(e))
	})
}

func Test_MaskFromPrefixLen(t *testing.T) {
	t.Run("MaskFromPrefixLen 1", func(t *testing.T) {
		m := MaskFromPrefixLen(24)
		e := &MACAddress{0xff, 0xff, 0xff, 0, 0, 0}
		assert.Equal(t, e, m)
	})
	t.Run("MaskFromPrefixLen 2", func(t *testing.T) {
		m := MaskFromPrefixLen(28)
		e := &MACAddress{0xff, 0xff, 0xff, 0xf0, 0, 0}
		assert.Equal(t, e, m)
	})
	t.Run("MaskFromPrefixLen handle too large", func(t *testing.T) {
		m := MaskFromPrefixLen(1000)
		e := &MACAddress{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
		assert.Equal(t, e, m)
	})
}

func Test_FromBytes(t *testing.T) {
	r := FromBytes(0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa)
	e := MustParseMACAddress("ff:ee:dd:cc:bb:aa")
	ne := MustParseMACAddress("01:23:45:67:89:ab")
	t.Run("FromBytes result equals MAC Address", func(t *testing.T) {
		assert.Equal(t, r, e)
	})
	t.Run("FromBytes result does not equal MAC Address", func(t *testing.T) {
		assert.NotEqual(t, r, ne)
	})
}

func Test_FromByteArray(t *testing.T) {
	r := FromByteArray([]byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa})
	e := MustParseMACAddress("ff:ee:dd:cc:bb:aa")
	ne := MustParseMACAddress("01:23:45:67:89:ab")
	t.Run("FromByteArray result equals MAC Address", func(t *testing.T) {
		assert.Equal(t, r, e)
	})
	t.Run("FromByteArray result does not equal MAC Address", func(t *testing.T) {
		assert.NotEqual(t, r, ne)
	})
}

func ExampleParseMACAddress() {
	mac, err := ParseMACAddress("00:00:5e:00:53:ab")
	if err != nil {
		panic(err)
	}
	fmt.Println(mac.String())
	// Output:
	// 00:00:5e:00:53:ab
}

func ExampleMustParseMACAddress() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.String())
	// Output:
	// 00:00:5e:00:53:ab
}

func ExampleFromBytes() {
	mac := FromBytes(0x00, 0x00, 0x5e, 0x00, 0x53, 0xab)
	fmt.Println(mac.String())
	// Output:
	// 00:00:5e:00:53:ab
}

func ExampleFromByteArray() {
	mac := FromByteArray([]byte{0x00, 0x00, 0x5e, 0x00, 0x53, 0xab})
	fmt.Println(mac.String())
	// Output:
	// 00:00:5e:00:53:ab
}

func ExampleMaskFromPrefixLen() {
	mask := MaskFromPrefixLen(24)
	fmt.Println(mask.String())
	// Output:
	// ff:ff:ff:00:00:00
}

func ExampleMACAddress_String() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.String())
	// Output:
	// 00:00:5e:00:53:ab
}

func ExampleMACAddress_Dots() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.Dots())
	// Output:
	// 0000.5e00.53ab
}

func ExampleMACAddress_Dashes() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.Dashes())
	// Output:
	// 00-00-5e-00-53-ab
}

func ExampleMACAddress_NoSeparators() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.NoSeparators())
	// Output:
	// 00005e0053ab
}

func ExampleMACAddress_Int() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.Int())
	// Output:
	// 1577079723
}

func ExampleMACAddress_Mask() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	mask := MustParseMACAddress("ff:ff:ff:ff:00:00")
	result := mac.Mask(mask)
	fmt.Println(result.String())
	// Output:
	// 00:00:5e:00:00:00
}

func ExampleMACAddress_Equal() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	other1 := MustParseMACAddress("00:00:5e:00:53:ab")
	other2 := MustParseMACAddress("00:00:5e:00:53:ba")
	fmt.Println(mac.Equal(other1))
	fmt.Println(mac.Equal(other2))
	// Output:
	// true
	// false
}

func ExampleMACAddress_Format() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	formatted := mac.Format("xx$$xx_-_xx@xx=xx.xx")
	fmt.Println(formatted)
	// Output:
	// 00$$00_-_5e@00=53.ab
}

func ExampleMACAddress_OUI() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	oui := mac.OUI()
	fmt.Println(oui)
	// Output:
	// 00:00:5e
}

func ExampleMACAddress_ByteString() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	byteString := mac.ByteString()
	fmt.Println(byteString)
	// Output:
	// {0,0,94,0,83,171}
}

func ExampleMACAddress_Next() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	next := mac.Next()
	fmt.Println(next.String())
	// Output:
	// 00:00:5e:00:53:ac
}

func ExampleMACAddress_Previous() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	prev := mac.Previous()
	fmt.Println(prev.String())
	// Output:
	// 00:00:5e:00:53:aa
}

func ExampleMACAddress_GreaterThan() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	other := MustParseMACAddress("00:00:5e:00:53:09")
	fmt.Println(mac.GreaterThan(other))
	fmt.Println(other.GreaterThan(mac))
	// Output:
	// true
	// false
}

func ExampleMACAddress_LessThan() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	other := MustParseMACAddress("00:00:5e:00:53:ff")
	fmt.Println(mac.LessThan(other))
	fmt.Println(other.LessThan(mac))
	// Output:
	// true
	// false
}

func ExampleMACAddress_GEqual() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	other := MustParseMACAddress("00:00:5e:00:53:09")
	fmt.Println(mac.GEqual(other))
	fmt.Println(mac.GEqual(mac))
	// Output:
	// true
	// true
}

func ExampleMACAddress_LEqual() {
	mac := MustParseMACAddress("00:00:5e:00:53:ab")
	other := MustParseMACAddress("00:00:5e:00:53:ff")
	fmt.Println(mac.LEqual(other))
	fmt.Println(mac.LEqual(mac))
	// Output:
	// true
	// true
}
