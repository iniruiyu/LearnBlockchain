package blockchain

type TxOutput struct {
	Value int64
	//脚本 公钥
	ScriptPublicKey string
}

func (txOutput *TxOutput) UnlockWithAddress(address string) bool {
	return txOutput.ScriptPublicKey == address
}
