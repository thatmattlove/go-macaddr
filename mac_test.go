package macaddr_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thatmattlove/go-macaddr"
	"github.com/thatmattlove/go-macaddr/internal/constant"
	"github.com/thatmattlove/go-macaddr/internal/read"
)

func TestMustParseMACAddress(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		assert.Panics(t, func() {
			macaddr.MustParseMACAddress("this should panic")
		})
	})
	t.Run("no panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			macaddr.MustParseMACAddress("01:23:45:67:89:ab")
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
		e := p[1]
		t.Run(fmt.Sprintf("parse %d", i+1), func(t *testing.T) {
			t.Parallel()
			r, err := macaddr.ParseMACAddress(p[0])
			require.NoError(t, err)
			assert.Equal(t, r.String(), e)
		})
	}
	errs := []string{
		"0123.4567.89az",
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}
	for i, s := range errs {
		t.Run(fmt.Sprintf("error %d", i+1), func(t *testing.T) {
			t.Parallel()
			_, err := macaddr.ParseMACAddress(s)
			require.Error(t, err)
		})
	}
}

func Test_MACAddress(t *testing.T) {
	s := "01:23:45:67:89:ab"
	m, err := macaddr.ParseMACAddress(s)
	assert.Nil(t, err)
	t.Run("Int() returns int", func(t *testing.T) {
		t.Parallel()
		i := m.Int()
		var e int64 = 1250999896491
		assert.Equal(t, e, i)
	})
	t.Run("Int() returns 0", func(t *testing.T) {
		t.Parallel()
		var m *macaddr.MACAddress
		i := m.Int()
		var e int64 = 0
		assert.Equal(t, e, i)
	})
	t.Run("Mask() works properly", func(t *testing.T) {
		t.Parallel()
		macOut := "01:23:45:00:00:00"
		maskIn := macaddr.MustParseMACAddress("ff:ff:ff:00:00:00")
		maskOut := m.Mask(maskIn)
		assert.Equal(t, macOut, maskOut.String())
	})
	t.Run("Mask() returns nil", func(t *testing.T) {
		t.Parallel()
		maskIn := macaddr.MACAddress{}
		maskOut := m.Mask(&maskIn)
		assert.Nil(t, maskOut)
	})
	t.Run("MACAddress.Equal()", func(t *testing.T) {
		t.Parallel()
		et := macaddr.MustParseMACAddress(s)
		ef := macaddr.MustParseMACAddress("12:34:56:78:9a:bc")
		assert.True(t, m.Equal(et))
		assert.False(t, m.Equal(ef))
	})
	t.Run("MACAddress.String()", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, m.String(), s)
	})
	t.Run("MACAddress.String() returns nil", func(t *testing.T) {
		t.Parallel()
		var m *macaddr.MACAddress
		assert.Equal(t, constant.NilStr, m.String())
	})
	t.Run("MACAddress.Dots()", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, m.Dots(), "0123.4567.89ab")
	})
	t.Run("MACAddress.Dashes()", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, m.Dashes(), "01-23-45-67-89-ab")
	})
	t.Run("MACAddress.NoSeparators()", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, m.NoSeparators(), "0123456789ab")
	})
	t.Run("MACAddress.OUI() 1", func(t *testing.T) {
		t.Parallel()
		oui := m.OUI()
		assert.Equal(t, "01:23:45", oui)
	})
	t.Run("MACAddress.OUI() 2", func(t *testing.T) {
		t.Parallel()
		oui := m.OUI(24)
		assert.Equal(t, "01:23:45", oui)
	})
	t.Run("MACAddress.OUI() 3", func(t *testing.T) {
		t.Parallel()
		m := macaddr.MustParseMACAddress("00:55:DA:80:01:23")
		oui := m.OUI(28)
		assert.Equal(t, "00:55:da:80:00:00/28", oui)
	})
	t.Run("MACAddress.OUI() nil", func(t *testing.T) {
		t.Parallel()
		var m *macaddr.MACAddress
		oui := m.OUI()
		assert.Equal(t, constant.NilStr, oui)
	})
	t.Run("MACAddress.ByteString() nil", func(t *testing.T) {
		t.Parallel()
		var m *macaddr.MACAddress
		assert.Equal(t, constant.NilStr, m.ByteString())
	})
	t.Run("MACAddress.ByteString() 1", func(t *testing.T) {
		t.Parallel()
		e := "{1,35,69,103,137,171}"
		assert.Equal(t, e, m.ByteString())
	})
	t.Run("MACAddress.Clone() nil", func(t *testing.T) {
		t.Parallel()
		var m *macaddr.MACAddress
		assert.Nil(t, m.Clone())
	})
	t.Run("MACAddress.Next() nil", func(t *testing.T) {
		t.Parallel()
		var m *macaddr.MACAddress
		assert.Nil(t, m.Next())
	})
	t.Run("MACAddress.Next()", func(t *testing.T) {
		t.Parallel()

		e := macaddr.MustParseMACAddress("01:23:45:67:89:ac")
		assert.Equal(t, e, m.Next())

		m := macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:ff")
		e = macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:ff")
		assert.Equal(t, e, m.Next())

		m = macaddr.MustParseMACAddress("ff:ff:ff:ff:fe:ff")
		e = macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:00")
		assert.Equal(t, e, m.Next())
	})
	t.Run("MACAddress.Previous()", func(t *testing.T) {
		t.Parallel()

		e := macaddr.MustParseMACAddress("01:23:45:67:89:aa")
		assert.Equal(t, e, m.Previous())

		m := macaddr.MustParseMACAddress("00:00:00:00:00:00")
		e = macaddr.MustParseMACAddress("00:00:00:00:00:00")
		assert.Equal(t, e, m.Previous())

		m = macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:01")
		e = macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:00")
		assert.Equal(t, e, m.Previous())
	})
	t.Run("MACAddress.GreaterThan() nil", func(t *testing.T) {
		t.Parallel()
		var m *macaddr.MACAddress
		mm := macaddr.FromBytes(0x01, 0x23, 0x45, 0x67, 0x89, 0xab)
		assert.False(t, m.GreaterThan(&macaddr.MACAddress{}))
		assert.False(t, mm.GreaterThan(nil))
	})
	t.Run("MACAddress.GreaterThan() 1", func(t *testing.T) {
		t.Parallel()
		e := &macaddr.MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaa}
		assert.True(t, m.GreaterThan(e))
	})
	t.Run("MACAddress.GreaterThan() 2", func(t *testing.T) {
		t.Parallel()
		e := &macaddr.MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaf}
		assert.False(t, m.GreaterThan(e))
	})
	// LessThan()
	t.Run("MACAddress.LessThan() nil", func(t *testing.T) {
		t.Parallel()
		var m *macaddr.MACAddress
		mm := macaddr.FromBytes(0x01, 0x23, 0x45, 0x67, 0x89, 0xab)
		assert.False(t, m.LessThan(&macaddr.MACAddress{}))
		assert.False(t, mm.LessThan(nil))
	})
	t.Run("MACAddress.LessThan() 1", func(t *testing.T) {
		t.Parallel()
		e := &macaddr.MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaf}
		assert.True(t, m.LessThan(e))
	})
	t.Run("MACAddress.LessThan() 2", func(t *testing.T) {
		t.Parallel()
		e := &macaddr.MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaa}
		assert.False(t, m.LessThan(e))
	})
	// GEqual()
	t.Run("MACAddress.GEqual() nil", func(t *testing.T) {
		t.Parallel()
		var m *macaddr.MACAddress
		mm := macaddr.FromBytes(0x01, 0x23, 0x45, 0x67, 0x89, 0xab)
		assert.False(t, m.GEqual(&macaddr.MACAddress{}))
		assert.False(t, mm.GEqual(nil))
	})
	t.Run("MACAddress.GEqual() 1", func(t *testing.T) {
		t.Parallel()
		e := &macaddr.MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaa}
		assert.True(t, m.GEqual(e))
	})
	t.Run("MACAddress.GEqual() 2", func(t *testing.T) {
		t.Parallel()
		e := &macaddr.MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaf}
		assert.False(t, m.GEqual(e))
	})
	t.Run("MACAddress.GEqual() 3", func(t *testing.T) {
		t.Parallel()
		e := &macaddr.MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xab}
		assert.True(t, m.GEqual(e))
	})
	// LEqual()
	t.Run("MACAddress.LEqual() nil", func(t *testing.T) {
		t.Parallel()
		var m *macaddr.MACAddress
		mm := macaddr.FromBytes(0x01, 0x23, 0x45, 0x67, 0x89, 0xab)
		assert.False(t, m.LEqual(&macaddr.MACAddress{}))
		assert.False(t, mm.LEqual(nil))
	})
	t.Run("MACAddress.LEqual() 1", func(t *testing.T) {
		t.Parallel()
		e := &macaddr.MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaf}
		assert.True(t, m.LEqual(e))
	})
	t.Run("MACAddress.LEqual() 2", func(t *testing.T) {
		t.Parallel()
		e := &macaddr.MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xaa}
		assert.False(t, m.LEqual(e))
	})
	t.Run("MACAddress.LEqual() 3", func(t *testing.T) {
		t.Parallel()
		e := &macaddr.MACAddress{0x01, 0x23, 0x45, 0x67, 0x89, 0xab}
		assert.True(t, m.LEqual(e))
	})
}

func Test_MaskFromPrefixLen(t *testing.T) {
	type pair struct {
		m *macaddr.MACAddress
		i int
	}
	pairs := []pair{
		{&macaddr.MACAddress{0xff, 0xff, 0xff, 0, 0, 0}, 24},
		{&macaddr.MACAddress{0xff, 0xff, 0xff, 0xf0, 0, 0}, 28},
		{&macaddr.MACAddress{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 1000},
		{&macaddr.MACAddress{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, 0},
	}
	for i, p := range pairs {
		p := p
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			r := macaddr.MaskFromPrefixLen(p.i)
			assert.Equal(t, p.m, r)
		})
	}
	for i := 0; i <= constant.MacBitLen; i++ {
		t.Run(fmt.Sprintf("len %d", i), func(t *testing.T) {
			t.Parallel()
			r := macaddr.MaskFromPrefixLen(i)
			l := read.PrefixLength(*r)
			assert.Equal(t, i, l)
		})
	}

}

func Test_FromBytes(t *testing.T) {
	r := macaddr.FromBytes(0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa)
	e := macaddr.MustParseMACAddress("ff:ee:dd:cc:bb:aa")
	ne := macaddr.MustParseMACAddress("01:23:45:67:89:ab")
	t.Run("FromBytes result equals MAC Address", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, r, e)
	})
	t.Run("FromBytes result does not equal MAC Address", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, r, ne)
	})
}

func Test_FromByteArray(t *testing.T) {
	r := macaddr.FromByteArray([]byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa})
	e := macaddr.MustParseMACAddress("ff:ee:dd:cc:bb:aa")
	ne := macaddr.MustParseMACAddress("01:23:45:67:89:ab")
	t.Run("FromByteArray result equals MAC Address", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, r, e)
	})
	t.Run("FromByteArray result does not equal MAC Address", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, r, ne)
	})
}

func ExampleParseMACAddress() {
	mac, err := macaddr.ParseMACAddress("00:00:5e:00:53:ab")
	if err != nil {
		panic(err)
	}
	fmt.Println(mac.String())
	// Output:
	// 00:00:5e:00:53:ab
}

func ExampleMustParseMACAddress() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.String())
	// Output:
	// 00:00:5e:00:53:ab
}

func ExampleFromBytes() {
	mac := macaddr.FromBytes(0x00, 0x00, 0x5e, 0x00, 0x53, 0xab)
	fmt.Println(mac.String())
	// Output:
	// 00:00:5e:00:53:ab
}

func ExampleFromByteArray() {
	mac := macaddr.FromByteArray([]byte{0x00, 0x00, 0x5e, 0x00, 0x53, 0xab})
	fmt.Println(mac.String())
	// Output:
	// 00:00:5e:00:53:ab
}

func ExampleMaskFromPrefixLen() {
	mask := macaddr.MaskFromPrefixLen(24)
	fmt.Println(mask.String())
	// Output:
	// ff:ff:ff:00:00:00
}

func ExampleMACAddress_String() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.String())
	// Output:
	// 00:00:5e:00:53:ab
}

func ExampleMACAddress_Dots() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.Dots())
	// Output:
	// 0000.5e00.53ab
}

func ExampleMACAddress_Dashes() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.Dashes())
	// Output:
	// 00-00-5e-00-53-ab
}

func ExampleMACAddress_NoSeparators() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.NoSeparators())
	// Output:
	// 00005e0053ab
}

func ExampleMACAddress_Int() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	fmt.Println(mac.Int())
	// Output:
	// 1577079723
}

func ExampleMACAddress_Mask() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	mask := macaddr.MustParseMACAddress("ff:ff:ff:ff:00:00")
	result := mac.Mask(mask)
	fmt.Println(result.String())
	// Output:
	// 00:00:5e:00:00:00
}

func ExampleMACAddress_Equal() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	other1 := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	other2 := macaddr.MustParseMACAddress("00:00:5e:00:53:ba")
	fmt.Println(mac.Equal(other1))
	fmt.Println(mac.Equal(other2))
	// Output:
	// true
	// false
}

func ExampleMACAddress_Format() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	formatted := mac.Format("xx$$xx_-_xx@xx=xx.xx")
	fmt.Println(formatted)
	// Output:
	// 00$$00_-_5e@00=53.ab
}

func ExampleMACAddress_OUI() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	oui := mac.OUI()
	fmt.Println(oui)
	// Output:
	// 00:00:5e
}

func ExampleMACAddress_ByteString() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	byteString := mac.ByteString()
	fmt.Println(byteString)
	// Output:
	// {0,0,94,0,83,171}
}

func ExampleMACAddress_Next() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	next := mac.Next()
	fmt.Println(next.String())
	// Output:
	// 00:00:5e:00:53:ac
}

func ExampleMACAddress_Previous() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	prev := mac.Previous()
	fmt.Println(prev.String())
	// Output:
	// 00:00:5e:00:53:aa
}

func ExampleMACAddress_GreaterThan() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	other := macaddr.MustParseMACAddress("00:00:5e:00:53:09")
	fmt.Println(mac.GreaterThan(other))
	fmt.Println(other.GreaterThan(mac))
	// Output:
	// true
	// false
}

func ExampleMACAddress_LessThan() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	other := macaddr.MustParseMACAddress("00:00:5e:00:53:ff")
	fmt.Println(mac.LessThan(other))
	fmt.Println(other.LessThan(mac))
	// Output:
	// true
	// false
}

func ExampleMACAddress_GEqual() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	other := macaddr.MustParseMACAddress("00:00:5e:00:53:09")
	fmt.Println(mac.GEqual(other))
	fmt.Println(mac.GEqual(mac))
	// Output:
	// true
	// true
}

func ExampleMACAddress_LEqual() {
	mac := macaddr.MustParseMACAddress("00:00:5e:00:53:ab")
	other := macaddr.MustParseMACAddress("00:00:5e:00:53:ff")
	fmt.Println(mac.LEqual(other))
	fmt.Println(mac.LEqual(mac))
	// Output:
	// true
	// true
}
