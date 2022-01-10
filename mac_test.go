package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MustParseMACAddr(t *testing.T) {
	t.Run("MustParseMACAddr should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseMACAddr("this should panic")
		})
	})
	t.Run("MustParseMACAddr should not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			// This should not panic.
			MustParseMACAddr("01:23:45:67:89:ab")
		})
	})
}

func Test_ParseMACAddr(t *testing.T) {
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
			r := MustParseMACAddr(p[0])
			e := p[1]
			assert.Equal(t, r.String(), e)
		})
	}
}

func Test_MACAddress(t *testing.T) {
	s := "01:23:45:67:89:ab"
	m, err := ParseMACAddr(s)
	assert.Nil(t, err)
	t.Run("Integer()", func(t *testing.T) {
		i := m.Integer()
		assert.Equal(t, i, 1250999896491)
	})
	t.Run("Mask()", func(t *testing.T) {
		macOut := "01:23:45:00:00:00"
		maskIn := MustParseMACAddr("ff:ff:ff:00:00:00")
		maskOut := m.Mask(maskIn)
		assert.Equal(t, macOut, maskOut.String())
	})
	t.Run("MACAddress.Equal()", func(t *testing.T) {
		et := MustParseMACAddr(s)
		ef := MustParseMACAddr("12:34:56:78:9a:bc")
		assert.True(t, m.Equal(et))
		assert.False(t, m.Equal(ef))
	})
	t.Run("MACAddress.String()", func(t *testing.T) {
		assert.Equal(t, m.String(), s)
	})
	t.Run("MACAddress.Dots()", func(t *testing.T) {
		assert.Equal(t, m.Dotted(), "0123.4567.89ab")
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
	e := MustParseMACAddr("ff:ee:dd:cc:bb:aa")
	ne := MustParseMACAddr("01:23:45:67:89:ab")
	t.Run("FromBytes result equals MAC Address", func(t *testing.T) {
		assert.Equal(t, r, e)
	})
	t.Run("FromBytes result does not equal MAC Address", func(t *testing.T) {
		assert.NotEqual(t, r, ne)
	})
}

func Test_FromByteArray(t *testing.T) {
	r := FromByteArray([6]byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa})
	e := MustParseMACAddr("ff:ee:dd:cc:bb:aa")
	ne := MustParseMACAddr("01:23:45:67:89:ab")
	t.Run("FromByteArray result equals MAC Address", func(t *testing.T) {
		assert.Equal(t, r, e)
	})
	t.Run("FromByteArray result does not equal MAC Address", func(t *testing.T) {
		assert.NotEqual(t, r, ne)
	})
}
