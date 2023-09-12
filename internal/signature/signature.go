package signature

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

// Sign converts any given value to JSON and signs it with the provided ecdsa private key
func Sign(v any, privateKey *ecdsa.PrivateKey) (*Signature, error) {
	m, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("could not marshal value to json: %w", err)
	}

	h := sha256.Sum256(m)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, h[:])
	if err != nil {
		return nil, fmt.Errorf("could not sign transaction: %w", err)
	}

	return &Signature{R: r, S: s}, nil
}
