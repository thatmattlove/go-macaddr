package read_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatmattlove/go-macaddr/internal/read"
)

func Test_PrefixLength(t *testing.T) {
	type pair struct {
		b []byte
		i int
	}
	tests := []pair{
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 48},
		{[]byte{0xff, 0xff, 0xff, 0x80, 0, 0}, 25},
		{[]byte{0xff, 0xff, 0xff, 0, 0, 0}, 24},
		{[]byte{0xff, 0xff, 0xff, 0xfa, 0, 0}, -1},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0, 0x5}, -1},
		{[]byte{0, 0, 0, 0, 0, 0}, 0},
	}
	for i := 0; i < len(tests); i++ {
		p := tests[i]
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			assert.Equal(t, p.i, read.PrefixLength(p.b))
		})
	}
}

func Test_IsZero(t *testing.T) {
	t.Parallel()
	assert.True(t, read.IsZero([]byte{0x0, 0x0, 0x0, 0x0}))
	assert.False(t, read.IsZero([]byte{0x1, 0x0, 0x0, 0x0}))
}

func Test_IsAllF(t *testing.T) {
	t.Parallel()
	assert.True(t, read.IsAllF([]byte{0xff, 0xff, 0xff}))
	assert.False(t, read.IsAllF([]byte{0xff, 0x00, 0xfe}))
}
