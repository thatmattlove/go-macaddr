package read

// PrefixLength is adapted from:
// https://github.com/golang/go/blob/2639f2f79bda2c3a4e9ef7381ca7de14935e2a4a/src/net/ip.go#L451
func PrefixLength(mac []byte) int {
	var n int
	for i, v := range mac {
		if v == 0xff {
			n += 8
			continue
		}
		// found non-ff byte
		// count 1 bits
		for v&0x80 != 0 {
			n++
			v <<= 1
		}
		// rest must be 0 bits
		if v != 0 {
			return -1
		}
		for i++; i < len(mac); i++ {
			if mac[i] != 0 {
				return -1
			}
		}
		break
	}
	return n
}
