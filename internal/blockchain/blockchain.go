package blockchain

import (
	"errors"
	"fmt"
	"strings"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
	address         string
}

func New(address string) (*Blockchain, error) {
	initialHash, err := (&Block{}).Hash()
	if err != nil {
		return nil, fmt.Errorf("could not get initial hash: %w", err)
	}

	bc := &Blockchain{address: address}
	bc.CreateBlock(0, initialHash)

	return bc, nil
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)

	return b
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value float32) {
	bc.transactionPool = append(bc.transactionPool, NewTransaction(sender, recipient, value))
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(t.SenderAddress, t.RecipientAddress, t.Value))
	}

	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) error {
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
	transactions := bc.CopyTransactionPool()
	previousHash, err := bc.LastBlock().Hash()
	if err != nil {
		return 0, err
	}

	var nonce int
	for bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) != nil {
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

	bc.AddTransaction(MINING_SENDER, bc.address, MINING_REWARD)
	bc.CreateBlock(nonce, previousHash)
	return nil
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var total float32
	for _, b := range bc.chain {
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
	return bc.chain[len(bc.chain)-1]
}
