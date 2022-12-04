package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"time"
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []Transaction
}

func NewBlock(nonce int, prevHash [32]byte, transactions []Transaction) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: prevHash,
		timestamp:    time.Now().UnixNano(),
		transactions: transactions,
	}
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int           `json:"nonce"`
		PreviousHash [32]byte      `json:"previous_hash"`
		Timestamp    int64         `json:"timestamp"`
		Transactions []Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}
