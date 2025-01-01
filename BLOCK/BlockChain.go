package block

import (
	"log"

	"github.com/boltdb/bolt"
)

// 数据库名字
const dbName = "blockchain.db"

// 表名字
const blockTableName = "block"

type Blockchain struct {
	Tip []byte   // 最新区块的hash
	DB  *bolt.DB // 数据库

}

//

/* type Blockchain struct {
	Blocks []*Block // 存储有序的区块
} */

// 增加区块到区块链当中
func (blockchain *Blockchain) AddBlockToBlockchain(data string) {
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		// 1. 获取表
		b := tx.Bucket([]byte(blockTableName))
		// 2. 创建新区块
		if b != nil {
			// 1. **先获取最新区块**
			TipBlockBytes := b.Get(blockchain.Tip)
			// 2. 反序列化
			TipBlock := DeSerializeBlock(TipBlockBytes)
			// 3. 将区块序列化并且存储到数据库中
			newBlock := NewBlock(data, TipBlock.Height+1, TipBlock.Hash)
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 存储最新的区块的hash
			err = b.Put([]byte("TipBlockHash"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})

}

// 创建区块链的方法

// 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain {
	var blockHash []byte
	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()  关闭不放在这里
	err = db.Update(func(tx *bolt.Tx) error {
		//创建数据库表
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}
		//表存在了
		if b != nil {
			// 创建 GenisisBlock
			GenesisBlock := CreateGenesisBlock("Genesis block Data...")
			// 将创世区块存储到表当中
			err := b.Put(GenesisBlock.Hash, GenesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 存储最新的区块的hash
			err = b.Put([]byte("TipBlockHash"), GenesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			blockHash = GenesisBlock.Hash
		}

		return nil
	})

	return &Blockchain{blockHash, db}
}
