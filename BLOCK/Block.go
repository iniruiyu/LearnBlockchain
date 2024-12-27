package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	// 1. 区块高度
	Height int64
	// 2. 上一个区块的hash
	PrevBlockHash []byte
	// 3. Data
	Data []byte
	// 4. 时间戳
	Timestamp int64
	// 5. Hash
	Hash []byte
}

// 2. 区块设置自己的hash
func (block *Block) SetHash() {

	// 1.将Height 转化为字节数组
	heightBytes := IntToHex(block.Height)
	// 2.将时间戳转换为字节数组
	// 2.1 将int64转换为字符串
	// 第二个参数范围为2~36，代表进制
	timeString := strconv.FormatInt(block.Timestamp, 2)
	timebytes := []byte(timeString)
	fmt.Println("SetHash timeString", timeString, "\n timebytes := ", timebytes, "\n heightBytes:=", heightBytes)
	// 3.拼接所有属性
	blockbytes := bytes.Join([][]byte{heightBytes, block.PrevBlockHash, block.Data, timebytes, block.Hash}, []byte{})
	// 4.生成hash
	HashValue := sha256.Sum256(blockbytes)

	// 处理 HashValue是32字节
	block.Hash = HashValue[:]
}

// 1. 创建新的区块
func NewBlock(data string, height int64, prevBlockHash []byte) *Block {
	// 创建区块
	block := &Block{
		Height:        height,
		PrevBlockHash: prevBlockHash,
		Data:          []byte(data),
		Timestamp:     time.Now().Unix(),
		Hash:          nil,
	}
	fmt.Println("old block = ", block)
	// 设置hash
	block.SetHash()
	return block
}