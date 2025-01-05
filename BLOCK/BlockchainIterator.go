package block

import (
	"log"

	"github.com/boltdb/bolt"
)

// 迭代器
type BlockchainIterator struct {
	CurrentHash []byte   // 最新区块的hash
	DB          *bolt.DB // 数据库
}

// 迭代器构造函数
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

// 迭代器方法
func (blockchainIterator *BlockchainIterator) NextPrevBlock() *Block {
	var currentBlock *Block
	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			currentBlockBytes := b.Get(blockchainIterator.CurrentHash)
			// 获取到当前迭代器Hash对应的区块
			currentBlock = DeSerializeBlock(currentBlockBytes)
			blockchainIterator.CurrentHash = currentBlock.PrevBlockHash

		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return currentBlock
}
