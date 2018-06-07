package blockchain

import (
	"fmt"
	"MiniGoBlockChain/src/cipher"
)

type Block struct {
	index uint64				// 区块号
	timeStamp uint64			// 时间戳
	previousHash cipher.SHA256  // 上一个区块的哈希值
	currentHash cipher.SHA256	// 当前区块的哈希值
}
func NewBlockChain(a string)  {
	fmt.Println("I am a block chain",a)
}

