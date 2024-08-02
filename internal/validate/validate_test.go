package validate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thatmattlove/go-macaddr/internal/constant"
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

func Test_ParseMacAddrWithPrefixLen(t *testing.T) {
	t.Run("parseMacAddrWithPrefixLen should error", func(t *testing.T) {
		m, p, e := validate.ParseMacAddrWithPrefixLen("this should error")
		assert.Empty(t, m)
		assert.Equal(t, 0, p)
		require.Error(t, e)
	})
	t.Run("parseMacAddrWithPrefixLen 2", func(t *testing.T) {
		m, p, e := validate.ParseMacAddrWithPrefixLen("01:23:45:67:89:ab/28")
		ms := "01:23:45:67:89:ab"
		assert.Equal(t, ms, m)
		assert.Equal(t, 28, p)
		require.NoError(t, e)
	})
	t.Run("parseMacAddrWithPrefixLen 3", func(t *testing.T) {
		m, p, e := validate.ParseMacAddrWithPrefixLen("00:00:00:00:00:00/0")
		ms := "00:00:00:00:00:00"
		assert.Equal(t, ms, m)
		assert.Equal(t, 0, p)
		require.NoError(t, e)
	})
	t.Run("parseMacAddrWithPrefixLen 4", func(t *testing.T) {
		m, p, e := validate.ParseMacAddrWithPrefixLen("01:23:45:67:89:ab")
		ms := "01:23:45:67:89:ab"
		assert.Equal(t, ms, m)
		assert.Equal(t, constant.MacBitLen, p)
		require.NoError(t, e)
	})
	t.Run("parseMacAddrWithPrefixLen 5", func(t *testing.T) {
		m, p, e := validate.ParseMacAddrWithPrefixLen("01:23:45:67:89:ab/ff")
		ms := "01:23:45:67:89:ab"
		assert.Equal(t, ms, m)
		assert.Equal(t, constant.MacBitLen, p)
		require.NoError(t, e)
	})
}
