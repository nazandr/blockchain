package blockchain

import (
	"bit-news/internal/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "BLOCKCHAIN"
	MINING_REWARD     = 10
)

type Blockchain struct {
	transactionPool []Transaction
	chain           []Block

	blockchainAddress string
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress

	bc.CreateBlock(0, b.Hash())

	return bc
}

func (bc *Blockchain) LastBlock() *Block {
	return &bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value *big.Int,
	senderPublicKey *ecdsa.PublicKey, signature *utils.Signature) error {
	t := NewTransaction(sender, recipient, value)

	if bc.CalculateTotalAmount(sender).Cmp(value) == -1 {
		return fmt.Errorf("%s has not enough balance at wallet", sender)
	}

	if sender != MINING_SENDER {
		v, err := bc.VerifyTransactionSignature(senderPublicKey, signature, t)
		if err != nil {
			return err
		}
		if !v {
			return fmt.Errorf("transaction signature not verified")
		}
	}

	bc.transactionPool = append(bc.transactionPool, *t)
	return nil
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, signature *utils.Signature,
	transaction *Transaction) (bool, error) {
	m, err := json.Marshal(transaction)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(m)

	return ecdsa.Verify(senderPublicKey, hash[:], signature.R, signature.S), nil
}

func (bc *Blockchain) CreateBlock(nonce int, prevHash [32]byte) *Block {
	b := NewBlock(nonce, prevHash, bc.transactionPool)
	bc.chain = append(bc.chain, *b)
	bc.transactionPool = []Transaction{}

	return b
}

func (bc *Blockchain) ValidProof(nonce int, prevHash [32]byte, transactions []Transaction, difficulty int) bool {
	guessBlock := Block{0, nonce, prevHash, transactions}

	return fmt.Sprintf("%x", guessBlock.Hash())[:difficulty] == strings.Repeat("0", difficulty)
}

func (bc *Blockchain) ProofOfWork(difficulty int) int {
	transactions := []Transaction{}
	copy(transactions, bc.transactionPool)
	prevHash := bc.LastBlock().Hash()

	nonce := 0
	for !bc.ValidProof(nonce, prevHash, transactions, difficulty) {
		nonce++
	}

	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, big.NewInt(MINING_REWARD), nil, nil)
	bc.CreateBlock(bc.ProofOfWork(MINING_DIFFICULTY), bc.LastBlock().Hash())
	log.Printf("[INFO] action=mining status=susses")

	return true
}

func (bc *Blockchain) CalculateTotalAmount(address string) *big.Int {
	totalAmount := big.NewInt(0)
	for _, block := range bc.chain {
		for _, transaction := range block.transactions {
			if address == transaction.recipientAddress {
				totalAmount = totalAmount.Add(totalAmount, transaction.value)
			}

			if address == transaction.senderAddress {
				totalAmount = totalAmount.Sub(totalAmount, transaction.value)
			}
		}
	}

	return totalAmount
}
