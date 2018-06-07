package main

import (
	"fmt"
	"MiniGoBlockChain/src/blockchain"
)

func main()  {
	blockchain.NewBlockChain("11111")
	blockchain.block()
	fmt.Println("Zibu's go block chain started")
}