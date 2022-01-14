package macaddr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
		t.Run(fmt.Sprintf("Parse MAC '%d'", i+1), func(t *testing.T) {
			r := MustParseMACAddress(p[0])
			e := p[1]
			assert.Equal(t, r.String(), e)
		})
	}
}

func Test_MACAddress(t *testing.T) {
	s := "01:23:45:67:89:ab"
	m, err := ParseMACAddress(s)
	assert.Nil(t, err)
	t.Run("Integer()", func(t *testing.T) {
		i := m.Int()
		assert.Equal(t, i, 1250999896491)
	})
	t.Run("Mask()", func(t *testing.T) {
		macOut := "01:23:45:00:00:00"
		maskIn := MustParseMACAddress("ff:ff:ff:00:00:00")
		maskOut := m.Mask(maskIn)
		assert.Equal(t, macOut, maskOut.String())
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
	t.Run("MACAddress.Dots()", func(t *testing.T) {
		assert.Equal(t, m.Dots(), "0123.4567.89ab")
	})
	t.Run("MACAddress.Dashes()", func(t *testing.T) {
		assert.Equal(t, m.Dashes(), "01-23-45-67-89-ab")
	})
	t.Run("MACAddress.NoSeparators()", func(t *testing.T) {
		assert.Equal(t, m.NoSeparators(), "0123456789ab")
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

func Test_createFmtString(t *testing.T) {
	t.Run("createFmtString 1", func(t *testing.T) {
		r := createFmtString("01:23:45:67:89:ab")
		assert.Equal(t, "xx:xx:xx:xx:xx:xx", r)
	})
	t.Run("createFmtString", func(t *testing.T) {
		r := createFmtString("0123.45:67-89ab")
		assert.Equal(t, "xxxx.xx:xx-xxxx", r)
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
