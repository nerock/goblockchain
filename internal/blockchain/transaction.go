package blockchain

type Transaction struct {
	SenderAddress    string  `json:"sender_blockchain_address"`
	RecipientAddress string  `json:"recipient_blockchain_address"`
	Value            float32 `json:"vlaue"`
}

func NewTransaction(sender, recipient string, value float32) *Transaction {
	return &Transaction{
		SenderAddress:    sender,
		RecipientAddress: recipient,
		Value:            value,
	}
}
