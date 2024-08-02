package convert_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatmattlove/go-macaddr/internal/constant"
	"github.com/thatmattlove/go-macaddr/internal/convert"
)

func Test_ChunkStr(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		t.Parallel()
		r := convert.ChunkStr("abcdefhij", 3)
		e := []string{"abc", "def", "hij"}
		assert.True(t, reflect.DeepEqual(r, e))
	})
	t.Run("returns empty 1", func(t *testing.T) {
		t.Parallel()
		r := convert.ChunkStr("", 0)
		assert.Len(t, r, 0)
	})
	t.Run("returns empty 2", func(t *testing.T) {
		t.Parallel()
		r := convert.ChunkStr("slice", 1)
		assert.Len(t, r, 5)
	})
}

func Test_DecToInt(t *testing.T) {
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
		p := p
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			t.Parallel()
			n, c, ok := convert.DecToInt(p.string)
			assert.True(t, ok)
			assert.Equal(t, n, p.int)
			assert.Equal(t, len(p.string), c)
		})
	}
	t.Run("Big", func(t *testing.T) {
		t.Parallel()
		n, c, ok := convert.DecToInt("16777215")
		assert.False(t, ok)
		assert.Equal(t, 7, c)
		assert.Equal(t, constant.Big, n)
	})
}

func Test_ByteArrayToInt64(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		t.Parallel()
		r := convert.ByteArrayToInt64([]byte{0xff})
		var e int64 = 255
		assert.Equal(t, e, r)
	})
}
