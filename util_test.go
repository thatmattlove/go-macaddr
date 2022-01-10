package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_chunkStr(t *testing.T) {
	t.Run("chunkStr works", func(t *testing.T) {
		r := chunkStr("abcdefhij", 3)
		e := []string{"abc", "def", "hij"}
		assert.True(t, reflect.DeepEqual(r, e))
	})
}

func Test_padRight(t *testing.T) {
	t.Run("padRight smaller", func(t *testing.T) {
		r := padRight("123", "0", 6)
		e := "123000"
		assert.True(t, reflect.DeepEqual(r, e))
	})
	t.Run("padRight larger", func(t *testing.T) {
		r := padRight("12345", "0", 3)
		e := "12345"
		assert.True(t, reflect.DeepEqual(r, e))
	})
}

func Test_reverseString(t *testing.T) {
	tests := [][]string{
		{"testing", "gnitset"},
		{"!@#$%^&*()", ")(*&^%$#@!"},
		{"🀱🀲🀳🀴🀵🀶🀷", "🀷🀶🀵🀴🀳🀲🀱"},
	}
	for i, p := range tests {
		t.Run(fmt.Sprintf("Pair %d", i), func(t *testing.T) {
			i := p[0]
			e := p[1]
			r := reverseString(i)
			assert.Equal(t, r, e)
		})
	}
}

func Test_decToInt(t *testing.T) {
	type pair = struct {
		string
		int
	}
	tests := []pair{
		{"123", 123},
		{"0123", 123},
		{"1024", 1024},
	}
	for i, p := range tests {
		t.Run(fmt.Sprintf("DecToInt %d", i), func(t *testing.T) {
			n, _, ok := decToInt(p.string)
			assert.True(t, ok)
			assert.Equal(t, n, p.int)
		})
	}
}

func Test_byteArrayToInt(t *testing.T) {
	t.Run("byteArrayToInt works", func(t *testing.T) {
		r := byteArrayToInt([]byte{0xff})
		assert.Equal(t, 255, r)
	})
}

func Test_validateHex(t *testing.T) {
	t.Run("validateMac false", func(t *testing.T) {
		f := validateHex("this is bs")
		assert.False(t, f)
	})
	t.Run("validateHex true", func(t *testing.T) {
		f := validateHex("01:23:45:67:89:ab")
		assert.True(t, f)
	})
}
