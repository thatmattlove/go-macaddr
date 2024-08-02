package constant

const (
	FmtDash  string = "xx-xx-xx-xx-xx-xx"
	FmtDot   string = "xxxx.xxxx.xxxx"
	FmtColon string = "xx:xx:xx:xx:xx:xx"
	FmtNone  string = "xxxxxxxxxxxx"
	NilStr   string = "<nil>"
)

const (
	Big                 int = 0xFFFFFF
	HexStrLen           int = 12
	HexStrWithColonsLen int = 17
	MacBitLen           int = 48
	MacByteLen          int = 6
)

var HexDigits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
