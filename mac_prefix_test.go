package macaddr

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatmattlove/go-macaddr/internal/constant"
)

func Test_MustParseMACPrefix(t *testing.T) {
	t.Run("MustParseMACPrefix should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseMACAddress("this should panic")
		})
	})
}
func Test_ParseMACPrefix(t *testing.T) {
	t.Run("Parse MAC Prefix", func(t *testing.T) {
		macE := "01:23:45:67:89:ab"
		maskE := "01:23:45:00:00:00/24"
		mac, mask, err := ParseMACPrefix("01:23:45:67:89:ab/24")
		assert.Nil(t, err)
		assert.Equal(t, mac.String(), macE)
		assert.Equal(t, mask.String(), maskE)
	})
	t.Run("ParseMACPrefix errors 1", func(t *testing.T) {
		_, _, err := ParseMACPrefix("this should error")
		assert.NotNil(t, err)
	})
	t.Run("ParseMACPrefix errors 2", func(t *testing.T) {
		m, mp, err := ParseMACPrefix("01:23:45:67:89:ab/64")
		assert.Nil(t, m)
		assert.Nil(t, mp)
		assert.NotNil(t, err)
	})
}

func Test_MACPrefix(t *testing.T) {
	s := "01:23:45:67:89:ab/24"
	_, mp, err := ParseMACPrefix(s)
	t.Run("Verify initial MACPrefix", func(t *testing.T) {
		assert.Nil(t, err)
		assert.NotNil(t, mp)
	})
	t.Run("MACPrefix.String()", func(t *testing.T) {
		assert.Equal(t, "01:23:45:00:00:00/24", mp.String())
	})
	t.Run("MACPrefix.String() is nil 1", func(t *testing.T) {
		var mp *MACPrefix
		assert.Equal(t, constant.NilStr, mp.String())
	})
	t.Run("MACPrefix.String() is nil 2", func(t *testing.T) {
		mp := MACPrefix{
			MAC:  nil,
			Mask: nil,
		}
		assert.Equal(t, constant.NilStr, mp.String())
	})
	t.Run("MACPrefix.Contains() 1", func(t *testing.T) {
		mc := MustParseMACAddress("01:23:45:ff:ee:dd")
		assert.True(t, mp.Contains(mc))
	})
	t.Run("MACPrefix.Contains() 2", func(t *testing.T) {
		_, mp := MustParseMACPrefix("44:6f:d8:10:00:00/28")
		mac := MustParseMACAddress("44:6f:d8:10:01:23")
		assert.True(t, mp.Contains(mac))
	})
	t.Run("MACPrefix.Contains() errors on nil prefix", func(t *testing.T) {
		mc := MustParseMACAddress("01:23:45:ff:ee:dd")
		var mp MACPrefix
		assert.Panics(t, func() {
			mp.Contains(mc)
		})
	})
	t.Run("MACPrefix.Contains() is false when lengths don't match", func(t *testing.T) {
		m := MACAddress{0xff, 0xff}
		assert.False(t, mp.Contains(&m))
	})
	t.Run("MACPrefix.PrefixLen()", func(t *testing.T) {
		type pair struct {
			string
			int
		}
		tests := []pair{
			{"01:23:45:67:89:ab/24", 24},
			{"01:23:45:00:00:00/48", 48},
			{"00:00:00:00:00:00/0", 0},
		}
		for _, p := range tests {
			_, r, _ := ParseMACPrefix(p.string)
			assert.Equal(t, r.PrefixLen(), p.int)
		}
	})
	t.Run("MACPrefix.Prefixlen() returns 0 when prefix is nil", func(t *testing.T) {
		var mp *MACPrefix
		assert.Equal(t, 0, mp.PrefixLen())
	})
	t.Run("MACPrefix.OUI() nil", func(t *testing.T) {
		var mp *MACPrefix
		assert.Equal(t, constant.NilStr, mp.OUI())
	})
	t.Run("MACPrefix.OUI() /24", func(t *testing.T) {
		_, mp := MustParseMACPrefix("01:23:45:00:00:00/24")
		assert.Equal(t, "01:23:45", mp.OUI())
	})
	t.Run("MACPrefix.OUI() /28", func(t *testing.T) {
		_, mp := MustParseMACPrefix("01:23:45:00:00:00/28")
		assert.Equal(t, "01:23:45:00:00:00/28", mp.OUI())
	})
	t.Run("MACPrefix.Match() matching MAC", func(t *testing.T) {
		m, e := mp.Match("01:23:45:67:89:ab")
		assert.Nil(t, e)
		assert.Equal(t, m.String(), mp.String())
	})
	t.Run("MACPrefix.Match() /28 matching MAC", func(t *testing.T) {
		_, mp := MustParseMACPrefix("44:6f:d8:10:00:00/28")
		m, err := mp.Match("44:6f:d8:10:01:23")
		assert.Nil(t, err)
		assert.Equal(t, m.String(), mp.String())
	})
	t.Run("MACPrefix.Match() matching OUI", func(t *testing.T) {
		m, e := mp.Match("01:23:45")
		assert.Nil(t, e)
		assert.Equal(t, m.String(), mp.String())
	})
	t.Run("MACPrefix.Match() non-matching", func(t *testing.T) {
		_, e := mp.Match("ba:98:76:54:32:01")
		assert.NotNil(t, e)
	})
	t.Run("MACPrefix.Match() non-matching OUI", func(t *testing.T) {
		_, e := mp.Match("98:76:54")
		assert.NotNil(t, e)
	})
	t.Run("MACPrefix.Match() larger prefix length", func(t *testing.T) {
		_, e := mp.Match("01:23:45:67:89:ab/12")
		assert.NotNil(t, e)
	})
	t.Run("MACPrefix.Match() invalid string", func(t *testing.T) {
		m, e := mp.Match("this should error")
		assert.Nil(t, m)
		assert.NotNil(t, e)
	})
	t.Run("MACPrefix.First() nil prefix", func(t *testing.T) {
		var mp *MACPrefix
		assert.Nil(t, mp.First())
	})
	t.Run("MACPrefix.First() 1", func(t *testing.T) {
		e := MustParseMACAddress("01:23:45:00:00:00")
		assert.Equal(t, e, mp.First())
	})
	t.Run("MACPrefix.Last() nil prefix", func(t *testing.T) {
		var mp *MACPrefix
		assert.Nil(t, mp.Last())
	})
	t.Run("MACPrefix.Last() 1", func(t *testing.T) {
		e := MustParseMACAddress("01:23:45:ff:ff:ff")
		assert.Equal(t, e, mp.Last())
	})
	t.Run("MACPrefix.Last() 2", func(t *testing.T) {
		_, mp := MustParseMACPrefix("00:11:22:30:00:00/28")
		e := MustParseMACAddress("00:11:22:3f:ff:ff")
		assert.Equal(t, e, mp.Last())
	})
	t.Run("MACPrefix.Count() nil prefix", func(t *testing.T) {
		var mp *MACPrefix
		assert.Equal(t, 0, mp.Count())
	})
	t.Run("MACPrefix.Count() /48", func(t *testing.T) {
		_, mp := MustParseMACPrefix("01:23:45:67:89:ab/48")
		assert.Equal(t, 1, mp.Count())
	})
	t.Run("MACPrefix.Count() /24", func(t *testing.T) {
		assert.Equal(t, 16_777_216, mp.Count())
	})
	t.Run("MACPrefix.Count() /28", func(t *testing.T) {
		_, mp := MustParseMACPrefix("01:23:45:67:89:ab/28")
		assert.Equal(t, 1_048_576, mp.Count())
	})
	t.Run("MACPrefix.Count() All MACs", func(t *testing.T) {
		_, mp := MustParseMACPrefix("00:00:00:00:00:00/0")
		e := int(math.Pow(2, 48))
		assert.Equal(t, e, mp.Count())
	})
	t.Run("MACPrefix.WildcardMask() nil prefix", func(t *testing.T) {
		var mp *MACPrefix
		assert.Nil(t, mp.WildcardMask())
	})
	t.Run("MACPrefix.WildcardMask() 1", func(t *testing.T) {
		e := MustParseMACAddress("00:00:00:ff:ff:ff")
		assert.Equal(t, e, mp.WildcardMask())
	})
	t.Run("MACPrefix.WildcardMask() 2", func(t *testing.T) {
		_, mp := MustParseMACPrefix("01:23:45:67:89:ab/28")
		e := MustParseMACAddress("00:00:00:0f:ff:ff")
		assert.Equal(t, e, mp.WildcardMask())
	})
	t.Run("MACPrefix.Iter() nil prefix", func(t *testing.T) {
		var mp *MACPrefix
		assert.Nil(t, mp.Iter())
	})
	t.Run("MACPrefix.Iter()", func(t *testing.T) {
		_, mp := MustParseMACPrefix("01:23:45:67:89:00/44")
		addrs := []MACAddress{
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x00},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x01},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x02},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x03},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x04},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x05},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x06},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x07},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x08},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x09},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x0a},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x0b},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x0c},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x0d},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x0e},
			{0x1, 0x23, 0x45, 0x67, 0x89, 0x0f},
		}
		iter := mp.Iter()
		for i := 0; iter.Next(); i++ {
			e := addrs[i]
			assert.Equal(t, &e, iter.Value())
		}
	})
}

func Test_parseMacAddrWithPrefixLen(t *testing.T) {
	t.Run("parseMacAddrWithPrefixLen should error", func(t *testing.T) {
		m, p, e := parseMacAddrWithPrefixLen("this should error")
		assert.Nil(t, m)
		assert.Equal(t, 0, p)
		assert.NotNil(t, e)
	})
	t.Run("parseMacAddrWithPrefixLen 2", func(t *testing.T) {
		m, p, e := parseMacAddrWithPrefixLen("01:23:45:67:89:ab/28")
		ms := "01:23:45:67:89:ab"
		assert.Equal(t, ms, m.String())
		assert.Equal(t, 28, p)
		assert.Nil(t, e)
	})
	t.Run("parseMacAddrWithPrefixLen 3", func(t *testing.T) {
		m, p, e := parseMacAddrWithPrefixLen("00:00:00:00:00:00/0")
		ms := "00:00:00:00:00:00"
		assert.Equal(t, ms, m.String())
		assert.Equal(t, 0, p)
		assert.Nil(t, e)
	})
	t.Run("parseMacAddrWithPrefixLen 4", func(t *testing.T) {
		m, p, e := parseMacAddrWithPrefixLen("01:23:45:67:89:ab")
		ms := "01:23:45:67:89:ab"
		assert.Equal(t, ms, m.String())
		assert.Equal(t, constant.MacBitLen, p)
		assert.Nil(t, e)
	})
}

func ExampleParseMACPrefix() {
	mac, macPrefix, err := ParseMACPrefix("00:00:5e:00:53:00/24")
	if err != nil {
		panic(err)
	}
	fmt.Println(mac.String())
	fmt.Println(macPrefix.String())
	fmt.Println(err)
	// Output:
	// 00:00:5e:00:53:00
	// 00:00:5e:00:00:00/24
	// <nil>
}

func ExampleMustParseMACPrefix() {
	mac, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	fmt.Println(mac.String())
	fmt.Println(macPrefix.String())
	// Output:
	// 00:00:5e:00:53:00
	// 00:00:5e:00:00:00/24
}

func ExampleMACPrefix_MAC() {
	_, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	fmt.Println(macPrefix.MAC)
	// Output:
	// 00:00:5e:00:00:00
}

func ExampleMACPrefix_Mask() {
	_, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	fmt.Println(macPrefix.Mask)
	// Output:
	// ff:ff:ff:00:00:00
}

func ExampleMACPrefix_Contains() {
	_, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	mac1 := MustParseMACAddress("00:00:5e:00:53:ab")
	mac2 := MustParseMACAddress("00:00:5f:00:53:ab")
	fmt.Println(macPrefix.Contains(mac1))
	fmt.Println(macPrefix.Contains(mac2))
	// Output:
	// true
	// false
}

func ExampleMACPrefix_Match() {
	mac, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	match, err := macPrefix.Match("00:00:5e")
	if err != nil {
		// Input string was not a match.
		panic(err)
	}
	fmt.Println(mac.String())
	fmt.Println(match.String())
	// Output:
	// 00:00:5e:00:53:00
	// 00:00:5e:00:00:00/24
}

func ExampleMACPrefix_OUI() {
	_, macPrefix1 := MustParseMACPrefix("00:00:5e:00:53:00/24")
	_, macPrefix2 := MustParseMACPrefix("00:00:5e:00:53:00/28")
	fmt.Println(macPrefix1.OUI())
	fmt.Println(macPrefix2.OUI())
	// Output:
	// 00:00:5e
	// 00:00:5e:00:00:00/28
}

func ExampleMACPrefix_PrefixLen() {
	_, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	fmt.Println(macPrefix.PrefixLen())
	// Output:
	// 24
}

func ExampleMACPrefix_String() {
	_, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	fmt.Println(macPrefix.String())
	// Output:
	// 00:00:5e:00:00:00/24
}

func ExampleMACPrefix_Count() {
	_, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	fmt.Println(macPrefix.Count())
	// Output:
	// 16777216
}

func ExampleMACPrefix_First() {
	_, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	fmt.Println(macPrefix.First())
	// Output:
	// 00:00:5e:00:00:00
}

func ExampleMACPrefix_Last() {
	_, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	fmt.Println(macPrefix.Last())
	// Output:
	// 00:00:5e:ff:ff:ff
}

func ExampleMACPrefix_WildcardMask() {
	_, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/24")
	fmt.Println(macPrefix.WildcardMask())
	// Output:
	// 00:00:00:ff:ff:ff
}

func ExampleMACPrefix_Iter() {
	_, macPrefix := MustParseMACPrefix("00:00:5e:00:53:00/44")
	iter := macPrefix.Iter()
	for iter.Next() {
		fmt.Println(iter.Value())
	}
	// Output:
	// 00:00:5e:00:53:00
	// 00:00:5e:00:53:01
	// 00:00:5e:00:53:02
	// 00:00:5e:00:53:03
	// 00:00:5e:00:53:04
	// 00:00:5e:00:53:05
	// 00:00:5e:00:53:06
	// 00:00:5e:00:53:07
	// 00:00:5e:00:53:08
	// 00:00:5e:00:53:09
	// 00:00:5e:00:53:0a
	// 00:00:5e:00:53:0b
	// 00:00:5e:00:53:0c
	// 00:00:5e:00:53:0d
	// 00:00:5e:00:53:0e
	// 00:00:5e:00:53:0f
}
