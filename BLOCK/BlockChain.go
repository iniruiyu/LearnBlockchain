package block

type Blockchain struct {
	Blocks []*Block // 存储有序的区块
}

// 增加区块到区块链当中

func (blockchain *Blockchain) AddBlockToBlockchain(data string, height int64, prevHash []byte) {
	newBlock := NewBlock(data, height, prevHash)
	// 往链数组上添加区块
	blockchain.Blocks = append(blockchain.Blocks, newBlock)
}

// 创建区块链的方法

// 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain {

	GenesisBlock := CreateGenesisBlock("Genesis block Data...")

	return &Blockchain{Blocks: []*Block{GenesisBlock}}
}
