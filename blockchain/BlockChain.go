package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

// 数据库名字
const dbName = "blockchain.db"

// 表名字
const blockTableName = "block"

// 数据库内最新hash名字
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
		fmt.Printf("-------------------------------Block Height %d---------------------------------------\n", block.Height)
		fmt.Printf("PrevBlockHash:%x ---", block.PrevBlockHash)
		fmt.Printf("Timestamp:%s \n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 PM"))
		fmt.Printf("Hash:%x --------------", block.Hash)
		fmt.Printf("Nonce:%d\n", block.Nonce)

		fmt.Printf("------------------------------Block %d----Txs:%v-------------------------\n", block.Height, block.Txs)
		for _, tx := range block.Txs {
			//fmt.Println(tx)
			fmt.Printf("Txs:  tx.hash=%x", tx.TxHash)
			fmt.Printf("\nVins →   ")
			for _, in := range tx.Vins {
				fmt.Println("in.TxHash:", hex.EncodeToString(in.TxHash))
				fmt.Println("数字签名:", in.ScriptSig)
				fmt.Printf("(in.Vout 索引): %d -----", in.Vout)
			}
			fmt.Printf("Vouts→   \n")
			for i, out := range tx.Vouts {
				fmt.Printf("tx.%d Value:%d-------", i, out.Value)
				fmt.Println("ScriptPublicKey:", out.ScriptPublicKey)
			}
			fmt.Printf("--------------------Block %d Ended-------------------------\n", block.Height)
		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

// 增加区块到区块链当中
func (blockchain *Blockchain) AddBlockToBlockchain(txs []*Transaction) {
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
			newBlock := NewBlock(txs, TipBlock.Height+1, TipBlock.Hash)
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
func CreateBlockchainWithGenesisBlock(genesisBlockAddress string) *Blockchain {

	//数据库存在直接退出
	if DbExists() {
		fmt.Println("创世区块已经存在")
		os.Exit(1)
		return nil
	}
	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 关闭数据库
	defer db.Close()

	var GenesisBlockHash []byte
	err = db.Update(func(tx *bolt.Tx) error {

		// 创建表
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}
		//表存在了
		if b != nil {
			// 创建 GenisisBlock
			// 创建了一个Coinbase Transaction
			txCoinbase := NewCoinbaseTransaction(genesisBlockAddress)
			GenesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
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

			GenesisBlockHash = GenesisBlock.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &Blockchain{GenesisBlockHash, db}
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

// 返回名下的所有未花费的交易输出
// 如果一个地址对应的TxOutput未花费，那么这个Transaction就应该添加到数组返回
func (blockchain *Blockchain) UnSpentTransactionsOutputWithAddress(address string, txs []*Transaction) []*UTXO {
	//var unSpentTxOutPut_UTXO []*TxOutput
	var unUTXOs []*UTXO
	//返回的这个Txoutput应该有一个对应的交易哈希
	//还对应一个Vout的索引
	// {hash:[0,1]}
	var spentTxOutput = make(map[string][]int)

	for _, tx := range txs {
		if !tx.IsCoinbaseTransaction() { //只有非创世区块才有Vin
			for _, in := range tx.Vins {
				// 是否能够解锁
				if in.UnlockWithAddress(address) {
					key := hex.EncodeToString(in.TxHash)

					spentTxOutput[key] = append(spentTxOutput[key], in.Vout)
					//fmt.Println("UnlockWithAddress VIN'S VOUT", in.Vout)
				}
			}
		}

	}

	//判断当前的有没有被消费
	for _, tx := range txs {

		for index, out := range tx.Vouts {
			if out.UnlockWithAddress(address) {
				if len(spentTxOutput) != 0 {
					var isSpendUTXO bool = false
					for hash, indexArray := range spentTxOutput {
						txhash := hex.EncodeToString(tx.TxHash)
						if txhash == hash {

							for _, outIndex := range indexArray {
								if index == outIndex {
									isSpendUTXO = true
									//continue work1
								}
							}
						} /*else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
							fmt.Println("2add", hex.EncodeToString(tx.TxHash))
						}*/
						if !isSpendUTXO {
							if out.Value > 0 {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						}
					}

				} else {
					if out.Value > 0 {
						utxo := &UTXO{tx.TxHash, index, out}
						unUTXOs = append(unUTXOs, utxo)
					}
				}
			}
		}
	}
	fmt.Println("unUTXO", unUTXOs)

	//通过迭代器进行遍历
	blockIterator := blockchain.Iterator()
	for {
		block := blockIterator.NextPrevBlock()
		//fmt.Println("block", block)

		//for _, tx := range block.Txs {
		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]
			// txhash
			// vin
			//把输入消费掉的存储到字典中
			if !tx.IsCoinbaseTransaction() { //只有非创世区块才有Vin
				for _, in := range tx.Vins {
					// 是否能够解锁
					if in.UnlockWithAddress(address) {
						key := hex.EncodeToString(in.TxHash)

						spentTxOutput[key] = append(spentTxOutput[key], in.Vout)
						//fmt.Println("UnlockWithAddress VIN'S VOUT", in.Vout)
					}
				}
			}

			//vout
			//work:
			for index, out := range tx.Vouts {
				//fmt.Println("out", out)
				// 是否能够解锁
				if spentTxOutput != nil {
					if out.UnlockWithAddress(address) {
						//fmt.Println("out", out)
						if len(spentTxOutput) != 0 {
							//判断是否被消费掉
							//var SpendTxHash []string
							var isSpendUTXO bool = false
							for txHash, indexArray := range spentTxOutput {

								for _, i := range indexArray {
									txHashStr := hex.EncodeToString(tx.TxHash)
									fmt.Println(txHashStr, txHash)
									if i == index && txHash == txHashStr {
										//说明当前字典在Vin的spentTxOutput字典中，已经花费掉了
										//SpendTxHash = append(SpendTxHash, txHash)
										fmt.Println("Address VIN'S VOUT Is Spended", out)
										isSpendUTXO = true
										//continue work
									}
								}
								//if !isSpendUTXO 在外面 error
							}
							if !isSpendUTXO {
								if out.Value > 0 {
									utxo := &UTXO{tx.TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
									fmt.Println("11add", hex.EncodeToString(tx.TxHash))
								}
							}
						} else {
							//unSpentTxOutPut_UTXO = append(unSpentTxOutPut_UTXO, out)
							fmt.Println("not in spentTxOutput", hex.EncodeToString(tx.TxHash))
							if out.Value > 0 {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
								fmt.Println("22add", hex.EncodeToString(tx.TxHash))
							}
						}

					}
				}
			}

		}
		//fmt.Println("UTXO", unSpentTxOutPut_UTXO)
		// 退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		//fmt.Println(block.PrevBlockHash, hashInt)
		/*
			    [0 0 175 80 210 83 117 43 116 91 221 148 123 3 137 170 42 235 7 146 23 193 156 164 120 116 127 122 62 95 140 32] {false [8679702544858975264 3092573893165161636 8384538761069037994 192761660929323]}
				[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] {false []}
		*/
		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			//Find Genesis Block
			break
		}
	}

	/*fmt.Println("unUTXOs", unUTXOs)
	for _, v := range unUTXOs {
		fmt.Println(v.Index)
		fmt.Println(hex.EncodeToString(v.TxHash))
		fmt.Println(v.Output)
	}
	Debug 重复
	*/

	return unUTXOs
}

func (blockchain *Blockchain) GetBalance(address string) int64 {
	utxos := blockchain.UnSpentTransactionsOutputWithAddress(address, []*Transaction{})
	var amount int64
	for i, utxo := range utxos {
		fmt.Println("utxo", i, "hash", hex.EncodeToString(utxo.TxHash), "Index", utxo.Index, "Output", utxo.Output)
		amount += utxo.Output.Value
	}
	return amount
}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendAbleUTXO(form string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	// 1.获取所有的UTXO
	UTXOs := blockchain.UnSpentTransactionsOutputWithAddress(form, txs)

	var SpendAbleUTXO map[string][]int = make(map[string][]int)
	// 2. 遍历utxos
	var value int64
	for _, utxo := range UTXOs {
		value = value + utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxHash)
		SpendAbleUTXO[hash] = append(SpendAbleUTXO[hash], utxo.Index)
		if value >= amount {
			break
		}
	}
	if value < amount {
		fmt.Println(form, "'s Fund is not enough")
		os.Exit(1)
	}
	return value, SpendAbleUTXO
}

// 挖掘新的区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {

	// 1. 建立一笔交易
	var txs []*Transaction
	// 2. 建立多笔交易
	for index, address := range from {
		amountValue, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], int64(amountValue), blockchain, txs)
		txs = append(txs, tx)
	}

	//fmt.Println("tx ===>", tx)
	//fmt.Println("挖掘新的区块")
	//fmt.Println("f->", from)
	//fmt.Println("t->", to)
	//fmt.Println("a->", amount)

	// 1. 通过相关算法建立transaction数组

	var block *Block
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte(tipBlockHash))  //获取最新区块hash 字节数组
			blockBytes := b.Get(hash)            //获取最新区块
			block = DeSerializeBlock(blockBytes) //反序列化获得区块对象
		}
		return nil
	})
	// 2. 建立新的区块
	block = NewBlock(txs, block.Height+1, block.Hash) //更新block为新区块
	// 3. 存储新区块数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			b.Put(block.Hash, block.Serialize())
			b.Put([]byte(tipBlockHash), block.Hash)
			blockchain.Tip = block.Hash
		}

		return nil
	})
}
