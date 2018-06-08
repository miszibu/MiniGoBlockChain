package main

import (
	"fmt"
	"../src/blockchain"
)

func main()  {
	geneBlock := &blockchain.Transaction{"1","2","3",4.2}
	fmt.Println(geneBlock)
	fmt.Println("Zibu's go block chain started")
}