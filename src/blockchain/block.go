package blockchain

type Block struct {
	Index int				// 区块号
	TimeStamp int64			    // 时间戳
	PreviousHash string         // 上一个区块的哈希值
	CurrentHash string	        // 当前区块的哈希值
	Nonce int 			    // 工作量证明，计算正确Hash值的次数
	Transactions []Transaction  // 区块所打包的交易
}
//var geneBlock Transaction
func NewBlock(){

}

