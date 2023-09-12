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
func Sign(v any, k *ecdsa.PrivateKey) (*Signature, error) {
	m, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("could not marshal value to json: %w", err)
	}

	h := sha256.Sum256(m)
	r, s, err := ecdsa.Sign(rand.Reader, k, h[:])
	if err != nil {
		return nil, fmt.Errorf("could not sign transaction: %w", err)
	}

	return &Signature{R: r, S: s}, nil
}

// Verify converts any given value to JSON and verifies its signature with the provided ecdsa public key
func (s *Signature) Verify(v any, k *ecdsa.PublicKey) error {
	m, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("could not marshal value to json: %w", err)
	}

	h := sha256.Sum256(m)
	if !ecdsa.Verify(k, h[:], s.R, s.S) {
		return fmt.Errorf("could not verify signature with transaction and public key: %w", err)
	}

	return nil
}
