package wallet

import (
	"crypto/ecdsa"
	"encoding/json"
)

type Transaction struct {
	senderKey        *ecdsa.PrivateKey
	senderAddress    string
	recipientAddress string
	value            float32
}

func NewTransaction(senderKey *ecdsa.PrivateKey, senderAddress, recipientAddress string, value float32) *Transaction {
	return &Transaction{
		senderKey:        senderKey,
		senderAddress:    senderAddress,
		recipientAddress: recipientAddress,
		value:            value,
	}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderAddress,
		Recipient: t.recipientAddress,
		Value:     t.value,
	})
}
