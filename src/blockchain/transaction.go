package blockchain

type Transaction struct {
	Id string			// 交易主键
	Sender string		// 发送方地址
	Recipient string    // 接收方地址
	Amount float64		// 交易金额
	IsValid bool		// 是否有效
	// 分布式数据库 怎么修改isvalid
	// 加入unspendpool概念
}

func (tx *Transaction)GetId() string {
	return tx.Id
}
