package blockchain

type Transaction struct {
	Id string			// 交易主键
	Sender string		// 发送方地址
	Recipient string    // 接收方地址
	Amount float64		// 交易金额
	IsValid bool		// 是否有效
}

func (tx *Transaction)GetId() string {
	return tx.Id
}
