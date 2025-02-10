package blockchain

type UTXO struct {
	TxHash []byte
	Index  int
	Output *TxOutput
}
