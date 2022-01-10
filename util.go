package macaddr

import (
	"regexp"
	"strings"
)

func byteArrayToInt(arr []byte) int {
	var res int
	for _, v := range arr {
		res <<= 8
		res |= int(v)
	}
	return res
}

func chunkStr(str string, size int) []string {
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

func padRight(str string, pad string, count int) string {
	n := count - len(str)
	if n > 0 {
		return str + strings.Repeat(pad, n)
	}
	return str
}

// DecToInt converts decimal string to integer.
// See Golang internal: https://github.com/golang/go/blob/a59e33224e42d60a97fa720a45e1b74eb6aaa3d0/src/net/parse.go#L122
func decToInt(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= BIG {
			return BIG, i, false
		}
	}
	if i == 0 {
		return 0, 0, false
	}
	return n, i, true
}

func reverseString(i string) (o string) {
	rr := []rune(i)
	var r []rune
	for i := len(rr) - 1; i >= 0; i-- {
		r = append(r, rr[i])
	}
	o = string(r)
	return
}

func validateHex(i string) (o bool) {
	p := regexp.MustCompile(`[a-zA-Z0-9]+`)
	hp := regexp.MustCompile(`[a-f0-9]+`)
	ss := p.FindAllString(i, -1)
	sj := strings.ToLower(strings.Join(ss, ""))
	hps := strings.Join(hp.FindAllString(sj, -1), "")
	return len(sj) == len(hps)
}
