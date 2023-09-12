package blockchain

import (
	"encoding/json"
	"fmt"
)

type Hash [32]byte

func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%x", h))
}
