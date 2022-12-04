package blockchain

import (
	"bit-news/internal/utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"math/big"
)

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            *big.Int
}

func NewTransaction(sender, recipient string, value *big.Int) *Transaction {
	return &Transaction{
		senderAddress:    sender,
		recipientAddress: recipient,
		value:            value,
	}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string   `json:"sender_address"`
		Recipient string   `json:"recipient_address"`
		Value     *big.Int `json:"value"`
	}{
		Sender:    t.senderAddress,
		Recipient: t.recipientAddress,
		Value:     t.value,
	})
}

type TransactionWithKeys struct {
	senderPrivateKey *ecdsa.PrivateKey
	senderPublicKey  *ecdsa.PublicKey

	senderAddress    string
	recipientAddress string
	value            *big.Int
}

func NewTransactionWithKeys(senderPrivateKey *ecdsa.PrivateKey, senderPublicKey *ecdsa.PublicKey,
	sender, recipient string, value *big.Int) *TransactionWithKeys {
	return &TransactionWithKeys{
		senderPrivateKey: senderPrivateKey,
		senderPublicKey:  senderPublicKey,
		senderAddress:    sender,
		recipientAddress: recipient,
		value:            value,
	}
}

func (t *TransactionWithKeys) GenerateSignature() (*utils.Signature, error) {
	m, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(m)

	r, s, err := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	if err != nil {
		return nil, err
	}

	return &utils.Signature{R: r, S: s}, nil
}

func (t *TransactionWithKeys) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string   `json:"sender_address"`
		Recipient string   `json:"recipient_address"`
		Value     *big.Int `json:"value"`
	}{
		Sender:    t.senderAddress,
		Recipient: t.recipientAddress,
		Value:     t.value,
	})
}
