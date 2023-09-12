package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nerock/goblockchain/internal/transaction"
)

type Block struct {
	Nonce        int                        `json:"nonce"`
	PreviousHash Hash                       `json:"previous_hash"`
	Timestamp    int64                      `json:"timestamp"`
	Transactions []*transaction.Transaction `json:"transactions"`
}

func NewBlock(nonce int, previousHash Hash, transactions []*transaction.Transaction) *Block {
	return &Block{
		Nonce:        nonce,
		PreviousHash: previousHash,
		Timestamp:    time.Now().UnixNano(),
		Transactions: transactions,
	}
}

func (b *Block) Hash() (Hash, error) {
	m, err := json.Marshal(b)
	if err != nil {
		return Hash{}, fmt.Errorf("could not marshal block to json: %w", err)
	}

	return sha256.Sum256(m), nil
}
