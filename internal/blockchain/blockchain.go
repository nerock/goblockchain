package blockchain

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"

	"github.com/nerock/goblockchain/internal/signature"
	"github.com/nerock/goblockchain/internal/transaction"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Blockchain struct {
	TransactionPool []*transaction.Transaction `json:"-"`
	Chain           []*Block                   `json:"blocks"`
	Address         string                     `json:"address"`
}

func New(address string) (*Blockchain, error) {
	initialHash, err := (&Block{}).Hash()
	if err != nil {
		return nil, fmt.Errorf("could not get initial hash: %w", err)
	}

	bc := &Blockchain{Address: address}
	bc.AddBlock(0, initialHash)

	return bc, nil
}

func (bc *Blockchain) AddBlock(nonce int, previousHash Hash) *Block {
	b := NewBlock(nonce, previousHash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, b)
	bc.TransactionPool = nil

	return b
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, sig *signature.Signature) error {
	t := transaction.New(sender, recipient, value)
	if sender != MINING_SENDER {
		if err := sig.Verify(t, senderPublicKey); err != nil {
			return err
		}
	}

	bc.TransactionPool = append(bc.TransactionPool, transaction.New(sender, recipient, value))
	return nil
}

func (bc *Blockchain) CopyTransactionPool() []*transaction.Transaction {
	transactions := make([]*transaction.Transaction, 0)
	for _, t := range bc.TransactionPool {
		transactions = append(transactions, transaction.New(t.SenderAddress, t.RecipientAddress, t.Value))
	}

	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash Hash, transactions []*transaction.Transaction, difficulty int) error {
	guessHash, err := NewBlock(nonce, previousHash, transactions).Hash()
	if err != nil {
		return err
	}
	guessHashStr := fmt.Sprintf("%x", guessHash)

	if guessHashStr[:difficulty] != strings.Repeat("0", difficulty) {
		return errors.New("invalid solution to challenge")
	}

	return nil
}

func (bc *Blockchain) ProofOfWork() (int, error) {
	if len(bc.TransactionPool) == 0 {
		return 0, errors.New("no transactions to mine")
	}

	previousHash, err := bc.LastBlock().Hash()
	if err != nil {
		return 0, err
	}

	var nonce int
	for bc.ValidProof(nonce, previousHash, bc.CopyTransactionPool(), MINING_DIFFICULTY) != nil {
		nonce++
	}

	return nonce, nil
}

func (bc *Blockchain) Mining() error {
	nonce, err := bc.ProofOfWork()
	if err != nil {
		return fmt.Errorf("generate proof of work: %w", err)
	}

	previousHash, err := bc.LastBlock().Hash()
	if err != nil {
		return fmt.Errorf("retrieve last block hash: %w", err)
	}

	if err := bc.AddTransaction(MINING_SENDER, bc.Address, MINING_REWARD, nil, nil); err != nil {
		return fmt.Errorf("adding reward transaction: %w", err)
	}

	bc.AddBlock(nonce, previousHash)
	return nil
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var total float32
	for _, b := range bc.Chain {
		for _, t := range b.Transactions {
			value := t.Value
			if blockchainAddress == t.SenderAddress {
				value *= -1
			}

			total += value
		}
	}

	return total
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}
