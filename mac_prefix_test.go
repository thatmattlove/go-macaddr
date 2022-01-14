package macaddr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
}

func Test_MACPrefix(t *testing.T) {
	s := "01:23:45:67:89:ab/24"
	_, mp, err := ParseMACPrefix(s)
	t.Run("Verify initial MACPrefix", func(t *testing.T) {
		assert.Nil(t, err)
		assert.NotNil(t, mp)
	})
	t.Run("MACPrefix.String()", func(t *testing.T) {
		assert.Equal(t, mp.String(), "01:23:45:00:00:00/24")
	})
	t.Run("MACPrefix.Contains()", func(t *testing.T) {
		mc := MustParseMACAddress("01:23:45:ff:ee:dd")
		assert.Nil(t, err)
		assert.True(t, mp.Contains(mc))
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
	t.Run("MACPrefix.OUI() /24", func(t *testing.T) {
		_, mp := MustParseMACPrefix("01:23:45:00:00:00/24")
		assert.Equal(t, "01:23:45", mp.OUI())
	})
	t.Run("MACPrefix.OUI() /28", func(t *testing.T) {
		_, mp := MustParseMACPrefix("01:23:45:00:00:00/28")
		assert.Equal(t, "01:23:45:00:00:00/28", mp.OUI())
	})
	t.Run("MACPrefix.Match() 1", func(t *testing.T) {
		m, e := mp.Match("01:23:45:67:89:ab")
		assert.Nil(t, e)
		assert.Equal(t, m.String(), mp.String())
	})
	t.Run("MACPrefix.Match() 2", func(t *testing.T) {
		m, e := mp.Match("01:23:45")
		assert.Nil(t, e)
		assert.Equal(t, m.String(), mp.String())
	})
	t.Run("MACPrefix.Match() 3", func(t *testing.T) {
		_, e := mp.Match("ba:98:76:54:32:01")
		assert.NotNil(t, e)
	})
	t.Run("MACPrefix.Match() 4", func(t *testing.T) {
		_, e := mp.Match("98:76:54")
		assert.NotNil(t, e)
	})
	t.Run("MACPrefix.Match() 5", func(t *testing.T) {
		_, e := mp.Match("01:23:45:67:89:ab/12")
		assert.NotNil(t, e)
	})
}

func Test_PrefixLength(t *testing.T) {
	type pair struct {
		MACAddress
		int
	}
	tests := []pair{
		{MACAddress{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 48},
		{MACAddress{0xff, 0xff, 0xff, 0, 0, 0}, 24},
		{MACAddress{0, 0, 0, 0, 0, 0}, 0},
	}
	for i, p := range tests {
		t.Run(fmt.Sprintf("prefixLength %d", i), func(t *testing.T) {
			assert.Equal(t, p.int, prefixLength(p.MACAddress))
		})
	}
}

func Test_parsePrefixLen(t *testing.T) {
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
		assert.Equal(t, 24, p)
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
