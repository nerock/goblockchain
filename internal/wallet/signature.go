package wallet

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

func Sign(t *Transaction) (*Signature, error) {
	m, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("could not conver transaction to json: %w", err)
	}

	h := sha256.Sum256(m)
	r, s, err := ecdsa.Sign(rand.Reader, t.senderKey, h[:])
	if err != nil {
		return nil, fmt.Errorf("could not sign transaction: %w", err)
	}

	return &Signature{R: r, S: s}, nil
}
