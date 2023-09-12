package blockchain

import "fmt"

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func New() (*Blockchain, error) {
	initialHash, err := (&Block{}).Hash()
	if err != nil {
		return nil, fmt.Errorf("could not get initial hash: %w", err)
	}

	bc := &Blockchain{}
	bc.CreateBlock(0, initialHash)

	return bc, nil
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)

	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}
