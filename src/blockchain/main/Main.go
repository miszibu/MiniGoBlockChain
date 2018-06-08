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
	mineBlock(sender)
	// 查询余额
	fmt.Println(sender,"地址的余额：",getWalletBalance(sender))
	// 发起一笔交易
	makeTransaction(sender,recipient,5)
	// 在挖一个矿
	mineBlock(sender)
	// 查询余额
	fmt.Println(sender,"地址的余额：",getWalletBalance(sender))
	// 显示unconfirmedpool
	fmt.Println(txs,len(txs))
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
func mineBlock(sender string)  {
	// 创建一个区块奖励交易 奖励10个币
	sysTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),"",sender,10,true}
	txs = append(txs, *sysTX)
	// 创建一个区块奖励交易 奖励1个币
	// syssTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),"",sender,1,true}
	//txs = append(txs, *syssTX)
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
	// 清空unconfirmpool
	//fmt.Println("挖矿后的区块链",txs,"  ",cap(txs))
	txs = make([]blockchain.Transaction,0)
	fmt.Println("挖矿后的区块链",BlockChain)
}

// 查询账户余额： 输入一个地址，遍历区块链中的所有交易，返回余额
func getWalletBalance(address string) float64{
	var balance float64
	for _,block := range BlockChain  {
		txs := block.Transactions
		for _,tx := range txs  {
			if address == tx.Recipient && tx.IsValid==true{
				balance += tx.Amount
			}
		}
	}
	return balance
}

// 发起一笔交易
func makeTransaction(sender string, recipient string, amount float64) {
	vaildTxs := make([]blockchain.Transaction,0)
	usedTxs := make([]blockchain.Transaction,0)
	transferValue := amount
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
	// 先交易零散金额 再交易大额UTXO
	// 选择validTx 加入 usedTx
	for _,tx := range  vaildTxs{
		if  tx.Amount>=transferValue{
			usedTxs = append(usedTxs, tx)
			transferValue -= tx.Amount
			break
		}else{
			transferValue -= tx.Amount
			usedTxs = append(usedTxs, tx)
		}
	}
	if transferValue>0 {
		fmt.Println("当前账户余额不足 无法完成交易")
		return
	}
	// 创建新的transaction
	for i,tx := range  usedTxs{
		// 前n-1个交易
		if i<cap(usedTxs)-1 {
			newTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),sender,recipient,tx.Amount,true}
			txs = append(txs, *newTX)
			amount-=tx.Amount
		}else { // 最后一个交易
			newTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),sender,recipient,amount,true}
			txs = append(txs, *newTX)
			changeTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),sender,sender,tx.Amount-amount,true}
			txs = append(txs, *changeTX)
		}
	}
	// 将使用过的transaction标记已使用
	for  i:=0;i<len(BlockChain);i++  {
		txs := BlockChain[i].Transactions
		for j:=0;j<len(txs);j++{
			for k:=0;k<len(usedTxs);k++{
				if txs[j].Id == usedTxs[k].Id {
					txs[j].IsValid=false
				}
			}
		}
	}
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
