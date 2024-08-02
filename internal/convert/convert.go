package convert

import "github.com/thatmattlove/go-macaddr/internal/constant"

// ByteArrayToInt converts a byte array to an int64.
func ByteArrayToInt64(arr []byte) int64 {
	var res int64
	for _, v := range arr {
		res <<= 8
		res |= int64(v)
	}
	return res
}

// ChunkStr chunks a string into chunks of n size. For example, "0123456789ab" with a size of 2
// would become [01 23 45 67 89 ab].
func ChunkStr(str string, size int) []string {
	var res []string
	slice := []rune(str)

	if len(slice) == 0 || size <= 0 {
		return res
	}

	length := len(slice)
	if size == 1 || size >= length {
		for _, v := range slice {
			var tmp rune
			tmp += v
			// tmp = append(tmp, v)
			res = append(res, string(tmp))
		}
		return res
	}

	// divide slice equally
	divideNum := length/size + 1
	for i := 0; i < divideNum; i++ {
		if i == divideNum-1 {
			if len(slice[i*size:]) > 0 {
				res = append(res, string(slice[i*size:]))
			}
		} else {
			res = append(res, string(slice[i*size:(i+1)*size]))
		}
	}

	return res
}

// DecToInt converts decimal string to integer.
// Returns number, characters consumed, success.
// See Golang internal: https://github.com/golang/go/blob/a59e33224e42d60a97fa720a45e1b74eb6aaa3d0/src/net/parse.go#L122
func DecToInt(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= constant.Big {
			return constant.Big, i, false
		}
	}
	if i == 0 {
		return 0, 0, false
	}
	return n, i, true
}
