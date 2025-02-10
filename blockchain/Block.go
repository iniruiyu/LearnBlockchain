package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Height        int64  // 1. 区块高度
	PrevBlockHash []byte // 2. 上一个区块的hash
	/* Data          []byte // 3. Data */
	Txs       []*Transaction // 3.交易数据
	Timestamp int64          // 4. 时间戳
	Hash      []byte         // 5. Hash
	Nonce     int64          // 6. Nonce （proof of work）
}

// 1. 创建新的区块
func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {
	// 创建区块
	block := &Block{
		Height:        height,
		PrevBlockHash: prevBlockHash,
		Txs:           txs,
		Timestamp:     time.Now().Unix(),
		Hash:          nil,
	}
	// 调用工作量证明方法，返回有效hash和Nonce值
	pow := NewProofOfWork(block)
	//000000
	// 挖矿验证
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

// 2. 单独写一个方法，生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	// 高度 hash可知
	var csqk [32]byte
	return NewBlock(txs, 1, csqk[:])
}

// 3. 返回交易哈希
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// 将区块序列化成字节数组
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化字节数组成区块
func DeSerializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
