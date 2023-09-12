package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

type Wallet struct {
	key     *ecdsa.PrivateKey
	address string
}

func New() (*Wallet, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("could not generate ecdsa key: %w", err)
	}

	return &Wallet{key: key, address: GenerateAddress(&key.PublicKey)}, nil
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.key
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return &w.key.PublicKey
}

func (w *Wallet) Address() string {
	return w.address
}
