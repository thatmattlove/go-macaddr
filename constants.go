package macaddr

const (
	_fmtDash  string = "xx-xx-xx-xx-xx-xx"
	_fmtDot   string = "xxxx.xxxx.xxxx"
	_fmtColon string = "xx:xx:xx:xx:xx:xx"
	_fmtNone  string = "xxxxxxxxxxxx"
	_nilStr   string = "<nil>"
)

const (
	_big                 int = 0xFFFFFF
	_hexStrLen           int = 12
	_hexStrWithColonsLen int = 17
	_macBitLen           int = 48
	_macByteLen          int = 6
)

var _hexDigits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
