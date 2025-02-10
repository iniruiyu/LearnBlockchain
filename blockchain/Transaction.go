package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// UTXO 未花费的交易输出
type Transaction struct {
	//1. 交易hash
	TxHash []byte
	//2. 输入
	Vins []*TxInput
	//3. 输出
	Vouts []*TxOutput
}

// 判断当前交易是否是Coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	//判断这笔交易是否是属于创世区块的那笔交易
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

// 将区块序列化成字节数组
func (tx *Transaction) HashTransactionSerialize() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

// Transaction 创建分两种情况
// 1. 创建创世区块时的Transaction
func NewCoinbaseTransaction(address string) *Transaction {

	// 代表消费
	txInput := &TxInput{[]byte{}, -1, "Genesis Data"}
	// 代表未消费(需要判断，曾经有这一笔钱)
	txOutput := &TxOutput{10, address}
	txCoinbase := &Transaction{[]byte{}, []*TxInput{txInput}, []*TxOutput{txOutput}}
	//设置hash值
	txCoinbase.HashTransactionSerialize()
	return txCoinbase
}

// 2. 转账时产生的Transaction

func NewSimpleTransaction(from string, to string, amount int64, blockchain *Blockchain, txs []*Transaction) *Transaction {
	//1. 创建一个函数，返回from地址的 所有未花费的交易输出所对应的Transaction

	//unSpentTransactionOutputUTXO := blockchain.UnSpentTransactionsOutputWithAddress(from)
	// 两种情况，同一个区块
	// {hash1,[0,1]}
	// 不同区块
	// {hash1,[0]，hash2,[1]}

	// 通过一个函数，返回我能花费的钱
	money, SpendAbleUTXO_Dic := blockchain.FindSpendAbleUTXO(from, amount, txs)

	//dic {hash,[0:2],hash,[1,4]}

	var txIntups []*TxInput
	var txOutputs []*TxOutput

	//代表消费
	//TxInput := &TxInput{[]byte("c2f026bb3e904c79e3be268e1fe04dcafec29f147f04dd3b896c15060c540140"), 0, from}
	//这么写有错
	for txHash, indexArray := range SpendAbleUTXO_Dic {

		bytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			//创建Input消费掉使用的Output
			TxInput := &TxInput{bytes, index, from}
			txIntups = append(txIntups, TxInput)

			str := hex.EncodeToString(bytes)
			fmt.Println("InputHash===========", str)
		}
	}

	//TxInput := &TxInput{bytes, 0, from}

	//消费
	// 代表未消费(需要判断，曾经有这一笔钱)

	//转账
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput = &TxOutput{money - amount, from}
	txOutputs = append(txOutputs, txOutput)

	txCoinbase := &Transaction{[]byte{}, txIntups, txOutputs}
	//设置hash值
	txCoinbase.HashTransactionSerialize()
	return txCoinbase
}
