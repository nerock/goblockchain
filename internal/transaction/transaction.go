package transaction

type Transaction struct {
	SenderAddress    string  `json:"sender_blockchain_address"`
	RecipientAddress string  `json:"recipient_blockchain_address"`
	Value            float32 `json:"value"`
}

func New(senderAddress, recipientAddress string, value float32) *Transaction {
	return &Transaction{
		SenderAddress:    senderAddress,
		RecipientAddress: recipientAddress,
		Value:            value,
	}
}
