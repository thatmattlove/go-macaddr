package validate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatmattlove/go-macaddr/internal/validate"
)

func Test_Hex(t *testing.T) {
	t.Run("validateMac false", func(t *testing.T) {
		f := validate.Hex("this is bs")
		assert.False(t, f)
	})
	t.Run("validateHex true", func(t *testing.T) {
		f := validate.Hex("01:23:45:67:89:ab")
		assert.True(t, f)
	})
}
