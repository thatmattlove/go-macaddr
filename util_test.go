package macaddr

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
	t.Run("chunkStr returns empty 1", func(t *testing.T) {
		r := chunkStr("", 0)
		assert.Len(t, r, 0)
	})
	t.Run("chunkStr returns empty 2", func(t *testing.T) {
		r := chunkStr("slice", 1)
		assert.Len(t, r, 5)
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

func Test_padMAC(t *testing.T) {
	t.Run("padMAC works", func(t *testing.T) {
		r := padMAC("012345")
		assert.Equal(t, "012345000000", r)
	})
	t.Run("padMAC invalids", func(t *testing.T) {
		r := padMAC("0x1y2z3")
		assert.Equal(t, "012300000000", r)
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

func Test_reverseBytes(t *testing.T) {
	t.Run("reverseBytes 1", func(t *testing.T) {
		b := []byte{0xff, 0xfe, 255, 1, 0}
		e := []byte{0, 1, 255, 0xfe, 0xff}
		assert.Equal(t, e, reverseBytes(b))
	})

	t.Run("reverseBytes empty", func(t *testing.T) {
		b := []byte{}
		e := []byte{}
		assert.Equal(t, e, reverseBytes(b))
	})
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
		t.Run(fmt.Sprintf("decToInt %d", i), func(t *testing.T) {
			n, c, ok := decToInt(p.string)
			assert.True(t, ok)
			assert.Equal(t, n, p.int)
			assert.Equal(t, len(p.string), c)
		})
	}
	t.Run("decToInt Big", func(t *testing.T) {
		n, c, ok := decToInt("16777215")
		assert.False(t, ok)
		assert.Equal(t, 7, c)
		assert.Equal(t, _big, n)
	})
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
