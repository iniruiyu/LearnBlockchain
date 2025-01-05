package block

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

// 数据库名字
const dbName = "blockchain.db"

// 表名字
const blockTableName = "block"

const tipBlockHash = "TipBlockHash"

type Blockchain struct {
	Tip []byte   // 最新区块的hash
	DB  *bolt.DB // 数据库
}

// 遍历输出所有区块的信息
func (blockchain *Blockchain) PrintBlockchain() {

	blockchainIterator := blockchain.Iterator()
	for {
		block := blockchainIterator.NextPrevBlock()

		fmt.Printf("Height:%d\n", block.Height)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Data:%s\n", block.Data)
		fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 PM"))
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Nonce:%d\n\n", block.Nonce)

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

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
			// 更新区块链Tip为最新的hash
			blockchain.Tip = newBlock.Hash
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 存储最新的区块的hash
			err = b.Put([]byte(tipBlockHash), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

}

// 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(genesisBlockData string) {
	/* var blockHash []byte

	var blockchain *Blockchain
	if dbExists() {
		// 创建或者打开数据库
		db, err := bolt.Open(dbName, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("GensisBlock alredy Existed")
		err = db.View(func(tx *bolt.Tx) error {
			// 获取表
			b := tx.Bucket([]byte(blockTableName))
			if b == nil {
				log.Panic(b)
			}
			tipHash := b.Get([]byte(tipBlockHash))
			if tipHash != nil {
				blockchain = &Blockchain{tipHash, db}
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

	} */
	//数据库存在直接退出
	if DbExists() {
		fmt.Println("创世区块已经存在")
		os.Exit(1)
		return
	}
	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//defer db.Close()  关闭不放在这里
	err = db.Update(func(tx *bolt.Tx) error {

		// 创建表
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}

		// 获取表
		/* 		b := tx.Bucket([]byte(blockTableName))
		   		if b == nil {
		   			//创建数据库表
		   			b, err = tx.CreateBucket([]byte(blockTableName))
		   			if err != nil {
		   				log.Panic(err)
		   			}
		   		} */
		//表存在了
		if b != nil {
			// 创建 GenisisBlock
			GenesisBlock := CreateGenesisBlock(genesisBlockData)
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
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// 判断数据库是否存在
func DbExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 获取blockchain对象
func GetBlockObject() *Blockchain {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	var TipHash []byte
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//最新区块哈希
			TipHash = b.Get([]byte("TipBlockHash"))
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return &Blockchain{TipHash, db}
}
