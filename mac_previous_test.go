package macaddr_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/go-macaddr"
)

// Test_Previous_BorrowPropagation covers the bug where Previous() of a non-zero
// MAC address whose least-significant byte(s) are 0x00 did not propagate the
// borrow across bytes, returning a value larger than the input (e.g.
// Previous("ff:ff:ff:ff:ff:00") returned "ff:ff:ff:ff:ff:ff" — the maximal
// address — instead of "ff:ff:ff:ff:fe:ff").
func Test_Previous_BorrowPropagation(t *testing.T) {
	t.Run("single-byte borrow", func(t *testing.T) {
		t.Parallel()
		cases := [][]string{
			{"00:00:5e:00:53:00", "00:00:5e:00:52:ff"},
			{"aa:bb:cc:dd:ee:00", "aa:bb:cc:dd:ed:ff"},
			{"01:23:45:67:89:00", "01:23:45:67:88:ff"},
			{"ff:ff:ff:ff:ff:00", "ff:ff:ff:ff:fe:ff"},
		}
		for _, c := range cases {
			c := c
			t.Run(c[0], func(t *testing.T) {
				t.Parallel()
				m := macaddr.MustParseMACAddress(c[0])
				e := macaddr.MustParseMACAddress(c[1])
				assert.Equal(t, e.String(), m.Previous().String())
			})
		}
	})
	t.Run("multi-byte borrow (two trailing zero bytes)", func(t *testing.T) {
		t.Parallel()
		m := macaddr.MustParseMACAddress("de:ad:be:ef:00:00")
		e := macaddr.MustParseMACAddress("de:ad:be:ee:ff:ff")
		assert.Equal(t, e.String(), m.Previous().String())
	})
	t.Run("MAX smoking gun: Previous of non-maximal address is NOT the maximal address", func(t *testing.T) {
		t.Parallel()
		m := macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:00")
		max := macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:ff")
		// Previous(x) must be strictly less than x; it must never return the
		// maximal address for a non-maximal input (an ordering inversion).
		assert.NotEqual(t, max.String(), m.Previous().String())
		assert.True(t, m.Previous().LessThan(m))
	})
	t.Run("Previous of MAX address", func(t *testing.T) {
		t.Parallel()
		m := macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:ff")
		e := macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:fe")
		assert.Equal(t, e.String(), m.Previous().String())
	})
}

// Test_Previous_SelfConsistency_Int verifies the library's own big-endian Int()
// model: Previous(x).Int() == x.Int() - 1 for every x with x.Int() > 0.
func Test_Previous_SelfConsistency_Int(t *testing.T) {
	t.Parallel()
	scan := []string{
		"00:00:5e:00:53:00",
		"aa:bb:cc:dd:ee:00",
		"ff:ff:ff:ff:ff:00",
		"de:ad:be:ef:00:00",
		"00:00:00:00:00:01",
		"00:00:00:00:00:10",
		"00:00:00:00:01:00",
		"00:00:00:01:00:00",
		"01:23:45:67:89:ab",
		"ff:ff:ff:ff:ff:ff",
	}
	for _, s := range scan {
		s := s
		t.Run(s, func(t *testing.T) {
			t.Parallel()
			m := macaddr.MustParseMACAddress(s)
			assert.Equal(t, m.Int()-1, m.Previous().Int())
		})
	}
}

// Test_Previous_InverseProperty verifies that Next and Previous are mutual
// inverses on (0, 2^48-1): Next(Previous(x)) == x and Previous(Next(x)) == x.
func Test_Previous_InverseProperty(t *testing.T) {
	t.Parallel()
	scan := []string{
		"00:00:5e:00:53:00",
		"aa:bb:cc:dd:ee:00",
		"ff:ff:ff:ff:ff:00",
		"de:ad:be:ef:00:00",
		"00:00:00:00:01:00",
		"00:00:00:01:00:00",
		"01:23:45:67:89:ab",
		"00:00:00:00:00:01",
	}
	for _, s := range scan {
		s := s
		t.Run(fmt.Sprintf("Next(Previous(%s))", s), func(t *testing.T) {
			t.Parallel()
			m := macaddr.MustParseMACAddress(s)
			require.Greater(t, m.Int(), int64(0))
			assert.Equal(t, s, m.Previous().Next().String())
		})
	}
	for _, s := range scan {
		s := s
		t.Run(fmt.Sprintf("Previous(Next(%s))", s), func(t *testing.T) {
			t.Parallel()
			m := macaddr.MustParseMACAddress(s)
			require.Less(t, m.Int(), int64(1<<48-1))
			assert.Equal(t, s, m.Next().Previous().String())
		})
	}
}

// Test_Previous_Controls verifies unchanged behavior for the no-borrow case and
// the zero-address clamp, plus Next carry propagation (the asymmetry signature).
func Test_Previous_Controls(t *testing.T) {
	t.Run("Previous no-borrow", func(t *testing.T) {
		t.Parallel()
		m := macaddr.MustParseMACAddress("00:00:00:00:00:01")
		e := macaddr.MustParseMACAddress("00:00:00:00:00:00")
		assert.Equal(t, e.String(), m.Previous().String())
	})
	t.Run("Previous zero-address clamp", func(t *testing.T) {
		t.Parallel()
		m := macaddr.MustParseMACAddress("00:00:00:00:00:00")
		e := macaddr.MustParseMACAddress("00:00:00:00:00:00")
		assert.Equal(t, e.String(), m.Previous().String())
	})
	t.Run("Next carry (single-byte)", func(t *testing.T) {
		t.Parallel()
		m := macaddr.MustParseMACAddress("ff:ff:ff:ff:fe:ff")
		e := macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:00")
		assert.Equal(t, e.String(), m.Next().String())
	})
	t.Run("Next carry (multi-byte)", func(t *testing.T) {
		t.Parallel()
		m := macaddr.MustParseMACAddress("de:ad:be:ef:ff:ff")
		e := macaddr.MustParseMACAddress("de:ad:be:f0:00:00")
		assert.Equal(t, e.String(), m.Next().String())
	})
	t.Run("Next all-f clamp", func(t *testing.T) {
		t.Parallel()
		m := macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:ff")
		e := macaddr.MustParseMACAddress("ff:ff:ff:ff:ff:ff")
		assert.Equal(t, e.String(), m.Next().String())
	})
}
