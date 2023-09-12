package cache

import (
	"sync"

	"github.com/nerock/goblockchain/internal/blockchain"
)

type BlockchainCache struct {
	m   sync.RWMutex
	bcs map[string]*blockchain.Blockchain
}

func New() *BlockchainCache {
	return &BlockchainCache{
		m:   sync.RWMutex{},
		bcs: make(map[string]*blockchain.Blockchain),
	}
}

func (c *BlockchainCache) Get(key string) *blockchain.Blockchain {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.bcs[key]
}

func (c *BlockchainCache) Put(key string, bc *blockchain.Blockchain) {
	c.m.Lock()
	defer c.m.Unlock()

	c.bcs[key] = bc
}
