package syntheticIP

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net"
	"testing"
)

func TestGenerateSyntheticIPAddress(t *testing.T) {
	t.Run("TestGenerateSyntheticIPAddress", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			_, cidr, err := net.ParseCIDR("2001:db8::/32")
			if err != nil {
				t.Fatalf("error in ParseCIDR: %s", err)
			}

			prefix := deepCopyCIDR(*cidr)
			seed := generateRandomString()
			ip, err := GenerateSyntheticIPAddress(seed, prefix)
			if err != nil {
				t.Fatalf("error in GenerateSyntheticIPAddress: %s", err)
			}
			if !prefix.Contains(ip) {
				t.Fatalf("generated IP address %s is not in prefix %s", ip, prefix)
			}
		}
	})
}

func generateRandomString() string {
	random := make([]byte, 32)
	io.ReadFull(rand.Reader, random)
	return base64.RawURLEncoding.EncodeToString(random)
}

func deepCopyCIDR(cidr net.IPNet) net.IPNet {
	ip := make([]byte, len(cidr.IP))
	copy(ip, cidr.IP)
	return net.IPNet{
		IP:   ip,
		Mask: cidr.Mask,
	}
}
