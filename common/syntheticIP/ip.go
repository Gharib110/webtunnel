package syntheticIP

import (
	"crypto/sha256"
	"net"
)

func GenerateSyntheticIPAddress(seed string, prefix net.IPNet) (net.IP, error) {
	hashValue := sha256.Sum256([]byte(seed))
	candidate := []byte(prefix.IP)
	one, _ := prefix.Mask.Size()
	for i := 0; i < len(candidate); i++ {
		if one < (i+1)*8 {
			candidate[i] = hashValue[i]
		} else {
			if one < i*8 {
				ones := i % 8
				candidate[i] ^= hashValue[i] & (0xff >> uint(8-ones))
			}
		}
	}
	return net.IP(candidate), nil
}
