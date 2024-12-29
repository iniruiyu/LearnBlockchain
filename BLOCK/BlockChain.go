package block

type Blockchain struct {
	Blocks []*Block // 存储有序的区块
}

// 创建区块链的方法

// 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain {

	GenesisBlock := CreateGenesisBlock("Genesis block Data...")

	return &Blockchain{Blocks: []*Block{GenesisBlock}}
}
