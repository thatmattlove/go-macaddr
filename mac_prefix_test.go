package macaddr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MustParseMACPrefix(t *testing.T) {
	t.Run("MustParseMACPrefix should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseMACAddr("this should panic")
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
	t.Run("MACPrefix.String()", func(t *testing.T) {
		assert.Equal(t, mp.String(), "01:23:45:00:00:00/24")
	})
	t.Run("MACPrefix.Contains()", func(t *testing.T) {
		mc := MustParseMACAddr("01:23:45:ff:ee:dd")
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
