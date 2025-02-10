package blockchain

type TxInput struct {
	// 1.交易哈希
	TxHash []byte

	// 2. 存储TxOutput在Vout里面的索引
	Vout int

	// 3. 用户名/数字签名
	ScriptSig string
}

// 判断当前消费的是谁的钱
func (TxInput *TxInput) UnlockWithAddress(address string) bool {
	return TxInput.ScriptSig == address
}
