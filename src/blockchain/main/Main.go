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
const INT_MAX = int(^uint(0) >> 1)

// 区块链
var BlockChain =  make([]blockchain.Block,0)
// unconfirmedpool
var unConfirmedTxs = make([]blockchain.Transaction,0)
// unspendpool
var unSpentTxs = make([]blockchain.Transaction,0)

func main(){
	// 初始化区块链：创建创世区块加入链中
	initBlockChain()
	// 挖一个区块，并获得区块奖励
	mineBlock(sender)
	// 查询余额
	fmt.Println(sender,"地址的余额：",getWalletBalance(sender))
	// 发起一笔交易
	makeTransaction(sender,recipient,5)
	// 再挖一个矿
	mineBlock(sender)
	// 查询余额
	fmt.Println(sender,"地址的余额：",getWalletBalance(sender))
	// 显示unconfirmedpool
	fmt.Println("unconfirmedpool",unConfirmedTxs,len(unConfirmedTxs))
	fmt.Println("unspentpool",unSpentTxs,len(unSpentTxs))

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
	sysTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),"",sender,10}
	unConfirmedTxs = append(unConfirmedTxs, *sysTX)
	// 创建一个区块奖励交易 奖励1个币
	// syssTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),"",sender,1,true}
	//unConfirmedTxs = append(unConfirmedTxs, *syssTX)
	// 获取当前区块链的最后一个区块
	latestBlock :=	BlockChain[cap(BlockChain)-1]
	var unConfirmedTxsString string
	if result,err:=json.Marshal(&unConfirmedTxs);err==nil{
		unConfirmedTxsString=string(result)
	}
	// 计算区块的解
	nonce := 1
	var hash string 
	for{
		hash = crtpto.SHA256(latestBlock.CurrentHash+unConfirmedTxsString+string(nonce))
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
		Transactions: unConfirmedTxs,
	}
	BlockChain = append(BlockChain,*newBlock)
	// 清空unconfirmpool
	unConfirmedTxs = make([]blockchain.Transaction,0)
	// 每当区块链延长时调用 更新unspentpool
	updateTxsToUnSpentPool();
	fmt.Println("当前区块链详情：",BlockChain)
}

//每当区块链延长时，读取最新的区块中的Txs，更新进入池中
func updateTxsToUnSpentPool(){
	txs := BlockChain[len(BlockChain)-1].Transactions
	for _,tx :=range txs {
		unSpentTxs=append(unSpentTxs,tx)
	}
}

//删除umspentpool中被消费的TX
func rmTxFormUnSpentPool(usedUnSpentTxs []blockchain.Transaction){
	for  i:=0;i<len(unSpentTxs);i++ {
		for j:=0;j<len(usedUnSpentTxs);j++ {
			if unSpentTxs[i].Id == usedUnSpentTxs[j].Id {
				unSpentTxs = append(unSpentTxs[:j], unSpentTxs[j+1:]...)
			}
		}
	}
}

// 查询账户余额： 查询unspendpool中的交易 返回余额
func getWalletBalance(address string) float64{
	var balance float64
	for _,tx := range unSpentTxs{
		if tx.Recipient == address {
			balance+=tx.Amount
		}
	}
	return balance
}

// 查询账户余额： 输入一个地址，遍历区块链中的所有交易，返回余额
/*func getWalletBalance(address string) float64{
	var balance float64
	for _,block := range BlockChain  {
		unConfirmedTxs := block.Transactions
		for _,tx := range unConfirmedTxs  {
			if address == tx.Recipient && tx.IsValid==true{
				balance += tx.Amount
			}
		}
	}
	return balance
}*/

// 发起一笔交易
func makeTransaction(sender string, recipient string, amount float64) {
	vaildUnSpentTxs := make([]blockchain.Transaction,0)
	usedUnSpentTxs := make([]blockchain.Transaction,0)
	transferValue := amount
	// 从unSpentPool中找到交易发起者的所有UTXO
	for _,tx := range unSpentTxs{
		if tx.Recipient == sender{
			vaildUnSpentTxs=append(vaildUnSpentTxs,tx)
		}
	}
	fmt.Println("found valid unSpentTxs",vaildUnSpentTxs)
	// 将需要validUnSpentTxs 按照Amount值 从小到大排序
	sort.Sort(txList{vaildUnSpentTxs})
	fmt.Println("found valid unSpentTxs",vaildUnSpentTxs)
	// 先交易零散金额 再交易大额UTXO
	// 选择validTx 加入 usedTx
	for _,tx := range  vaildUnSpentTxs{
		if  tx.Amount>=transferValue{
			usedUnSpentTxs = append(usedUnSpentTxs, tx)
			transferValue -= tx.Amount
			break
		}else{
			transferValue -= tx.Amount
			usedUnSpentTxs = append(usedUnSpentTxs, tx)
		}
	}
	if transferValue>0 {
		fmt.Println("当前账户余额不足 无法完成交易")
		return
	}
	// 创建新的transaction
	for i,tx := range  usedUnSpentTxs{
		// 前n-1个交易
		if i<cap(usedUnSpentTxs)-1 {
			newTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),sender,recipient,tx.Amount}
			unConfirmedTxs = append(unConfirmedTxs, *newTX)
			amount-=tx.Amount
		}else { // 最后一个交易
			newTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),sender,recipient,amount}
			unConfirmedTxs = append(unConfirmedTxs, *newTX)
			changeTX := &blockchain.Transaction{crtpto.RandGeneratorString(INT_MAX),sender,sender,tx.Amount-amount}
			unConfirmedTxs = append(unConfirmedTxs, *changeTX)
		}
	}
	// 从unspentpool中移除使用过的transaction
	rmTxFormUnSpentPool(usedUnSpentTxs)
}

// 调用 sort.Sort 需要重写三个方法
type txList struct {
	unConfirmedTxs []blockchain.Transaction
}
func (tx txList) Less(i, j int) bool {
	return tx.unConfirmedTxs[i].Amount < tx.unConfirmedTxs[j].Amount
}
func (tx txList) Len() int {
	return len(tx.unConfirmedTxs)
}
func (tx txList) Swap(i, j int) {
	tx.unConfirmedTxs[i], tx.unConfirmedTxs[j] = tx.unConfirmedTxs[j], tx.unConfirmedTxs[i]
}
