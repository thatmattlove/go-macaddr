package format_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mdl.wtf/go-macaddr/internal/format"
)

func Test_PadRight(t *testing.T) {
	t.Run("smaller", func(t *testing.T) {
		t.Parallel()
		r := format.PadRight("123", "0", 6)
		e := "123000"
		assert.True(t, reflect.DeepEqual(r, e))
	})
	t.Run("larger", func(t *testing.T) {
		t.Parallel()
		r := format.PadRight("12345", "0", 3)
		e := "12345"
		assert.True(t, reflect.DeepEqual(r, e))
	})

}

func Test_PadMAC(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		t.Parallel()
		r := format.PadMAC("012345")
		assert.Equal(t, "012345000000", r)
	})
	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		r := format.PadMAC("0x1y2z3")
		assert.Equal(t, "012300000000", r)
	})
}

func Test_ReverseString(t *testing.T) {
	tests := [][]string{
		{"testing", "gnitset"},
		{"!@#$%^&*()", ")(*&^%$#@!"},
		{"🀱🀲🀳🀴🀵🀶🀷", "🀷🀶🀵🀴🀳🀲🀱"},
	}
	for i, p := range tests {
		p := p
		t.Run(fmt.Sprintf("Pair %d", i+1), func(t *testing.T) {
			t.Parallel()
			i := p[0]
			e := p[1]
			r := format.ReverseString(i)
			assert.Equal(t, r, e)
		})
	}
}

func Test_ReverseBytes(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		t.Parallel()
		b := []byte{0xff, 0xfe, 255, 1, 0}
		e := []byte{0, 1, 255, 0xfe, 0xff}
		assert.Equal(t, e, format.ReverseBytes(b))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		b := []byte{}
		e := []byte{}
		assert.Equal(t, e, format.ReverseBytes(b))
	})
}

func Test_CreateFmtString(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		t.Parallel()
		r := format.CreateFmtString("01:23:45:67:89:ab")
		assert.Equal(t, "xx:xx:xx:xx:xx:xx", r)
	})
	t.Run("2", func(t *testing.T) {
		t.Parallel()
		r := format.CreateFmtString("0123.45:67-89ab")
		assert.Equal(t, "xxxx.xx:xx-xxxx", r)
	})
}

func Test_WithColons(t *testing.T) {
	t.Parallel()
	result := format.WithColons("0123456789ab")
	assert.Equal(t, "01:23:45:67:89:ab", result)
}
