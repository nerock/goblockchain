package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// GenerateAddress generates a blockchain wallet address based on an ecdsa public key
func GenerateAddress(publicKey *ecdsa.PublicKey) string {
	// SHA-256 hashing on the wallet public key
	h2 := sha256.New()
	h2.Write(publicKey.X.Bytes())
	h2.Write(publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	// RIPEMD-160 on the SHA-256 hash
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	// Add version byte in front of the RIPEMD-160 hash
	vd4 := make([]byte, 21)
	vd4[0] = 0x0
	copy(vd4[1:], digest3)

	// SHA-256 on the byte prepended RIPEMD-160 hash
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	// SHA-256 on the last SHA-256 hash
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	// Take the 4 first bytes of the last SHA-256 hash as checksum
	checksum := digest6[:4]

	// Join the checksum and the byte prepended RIPEMD-160 hash
	dc8 := make([]byte, 25)
	copy(dc8[21:], checksum)
	copy(dc8[:21], vd4)

	// Encode result to base58 as the wallet address
	return base58.Encode(dc8)
}
