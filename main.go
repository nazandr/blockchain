package main

import (
	"bit-news/internal/blockchain"
	"bit-news/internal/wallet"
	"fmt"
	"log"
	"math/big"
)

func main() {
	walletM, _ := wallet.NewWallet()
	walletA, _ := wallet.NewWallet()
	walletB, _ := wallet.NewWallet()

	bc := blockchain.NewBlockchain(walletM.BlockchainAddress())
	fmt.Printf("%+v\n\n", bc)

	t := blockchain.NewTransactionWithKeys(walletA.PrivateKey(), walletA.PublicKey(),
		walletA.BlockchainAddress(), walletB.BlockchainAddress(), big.NewInt((1)))
	tSig, _ := t.GenerateSignature()

	if err := bc.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), big.NewInt(1), walletA.PublicKey(),
		tSig); err != nil {
		log.Println(err)
	} else {
		log.Println("transaction added")
	}

	bc.Mining()
	fmt.Printf("%+v\n\n", bc)

	fmt.Println("A: ", bc.CalculateTotalAmount(walletA.BlockchainAddress()).String())
	fmt.Println("B: ", bc.CalculateTotalAmount(walletB.BlockchainAddress()).String())
	fmt.Println("M: ", bc.CalculateTotalAmount(walletM.BlockchainAddress()).String())
}
