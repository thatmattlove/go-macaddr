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
		{[]byte{0xff, 0xff, 0xff, 0, 0, 0}, 24},
		{[]byte{0, 0, 0, 0, 0, 0}, 0},
	}
	for i := 0; i < len(tests); i++ {
		p := tests[i]
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			assert.Equal(t, p.i, read.PrefixLength(p.b))
		})
	}

}
