package main

import (
	"fmt"
	"MiniGoBlockChain/src/blockchain"
	"MiniGoBlockChain/src/crtpto"
	"time"
	"encoding/json"
	"strings"
	"sort"
)

const sender   = "0x123456123456"
const recipient  = "0x111111111111"
const UINT_MAX = ^uint(0)
const INT_MAX = int(^uint(0) >> 1)

// 区块链
var BlockChain =  make([]blockchain.Block,0)
// unconfirmedpool
var txs = make([]blockchain.Transaction,0)


func main(){
	// 创建创始区块 并加入区块链里中
	initBlockChain()
	// 挖一个区块，并获得区块奖励
	mineBlock(txs, sender)
	// 查询余额
	fmt.Println(sender,"地址的余额：",getWalletBalance(sender))
	makeTransaction(sender,recipient,5)
}

// 创建创始区块 并加入区块链里中
func initBlockChain()  {
	// 生成创世区块
	genenisBlock := &blockchain.Block{
		Index:        0,
		TimeStamp:    time.Now().Unix(),
		PreviousHash: "1",
		CurrentHash:  "currentHash",
		Nonce:        1,
		Transactions: nil,
	}
	// 将创世区块加入区块链
	BlockChain= append(BlockChain, *genenisBlock)
}

 // 传入区块链 需记录的转账记录 挖矿
func mineBlock(txs []blockchain.Transaction,sender string)  {
	// 创建一个区块奖励交易 奖励10个币
	sysTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),"",sender,10,true}
	txs = append(txs, *sysTX)
	// 创建一个区块奖励交易 奖励10个币
	syssTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),"",sender,1,true}
	txs = append(txs, *syssTX)
	// 获取当前区块链的最后一个区块
	latestBlock :=	BlockChain[cap(BlockChain)-1]
	var txsString string
	if result,err:=json.Marshal(&txs);err==nil{
		txsString=string(result)
	}
	// 计算区块的解
	nonce := 1
	var hash string 
	for{
		hash = crtpto.SHA256(latestBlock.CurrentHash+txsString+string(nonce))
		if strings.HasPrefix(hash,"0000") {
			fmt.Println("计算结果正确 我获得了一个区块. nonce:",nonce,"hash:",hash)
			break
		}
		nonce++
		fmt.Println("第",nonce,"次计算错误，错误HASH为",hash)
	}
	// 解出结果 构造新区块，加入区块链中
	newBlock := &blockchain.Block{
		Index:        cap(BlockChain),
		TimeStamp:    time.Now().Unix(),
		PreviousHash: BlockChain[cap(BlockChain)-1].CurrentHash,
		CurrentHash:  hash,
		Nonce:        nonce,
		Transactions: txs,
	}
	BlockChain = append(BlockChain,*newBlock)
	fmt.Println("挖矿后的区块链",BlockChain)
	//return BlockChain
}

// 查询账户余额： 输入一个地址，遍历区块链中的所有交易，返回余额
func getWalletBalance(address string) float64{
	var balance float64
	for _,block := range BlockChain  {
		txs := block.Transactions
		for _,tx := range txs  {
			if address == tx.Recipient {
				balance += tx.Amount
			}
			if address == tx.Sender {
				balance -= tx.Amount
			}
		}
	}
	return balance
}

// 发起一笔交易
func makeTransaction(sender string, recipient string, amount float64) bool {
	isSuccess := true
	vaildTxs := make([]blockchain.Transaction,0)
	usedTxs := make([]blockchain.Transaction,0)
	//先找到一笔足够交易的钱
	for _,block := range BlockChain  {
		txs := block.Transactions
		for _,tx := range txs  {
			if tx.IsValid == true && tx.Recipient == sender {
				vaildTxs=append(vaildTxs,tx)
			}
		}
	}
	fmt.Println("found valid txs",vaildTxs)
	// 将需要validTxs 按照Amount值 从小到大排序
	sort.Sort(txList{vaildTxs})
	fmt.Println("found valid txs",vaildTxs)
	//
	for _,tx := range  vaildTxs{
		if  tx.Amount>=amount{
			usedTxs = append(usedTxs, tx)
			break
		}else{
			amount -= tx.Amount
			usedTxs = append(usedTxs, tx)
		}
	}
	if amount>0 {
		return  false
	}
	return isSuccess
}

// 调用 sort.Sort 需要重写三个方法
type txList struct {
	txs []blockchain.Transaction
}
func (tx txList) Less(i, j int) bool {
	return tx.txs[i].Amount < tx.txs[j].Amount
}
func (tx txList) Len() int {
	return len(tx.txs)
}
func (tx txList) Swap(i, j int) {
	tx.txs[i], tx.txs[j] = tx.txs[j], tx.txs[i]
}
