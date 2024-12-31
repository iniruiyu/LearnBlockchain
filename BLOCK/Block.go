package block

import (
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

	// 6.nonce （proof of work）
	Nonce int64
}

// 2. 区块设置自己的hash
/* func (block *Block) SetHash() {

	// 1.将Height 转化为字节数组
	heightBytes := IntToHex(block.Height)
	// 2.将时间戳转换为字节数组
	// 2.1 将int64转换为字符串
	// 第二个参数范围为2~36，代表进制
	timeString := strconv.FormatInt(block.Timestamp, 2)
	timebytes := []byte(timeString)
	//fmt.Println("SetHash timeString", timeString, "----- timebytes := ", timebytes, "------ heightBytes:=", heightBytes)
	// 3.拼接所有属性
	blockbytes := bytes.Join([][]byte{heightBytes, block.PrevBlockHash, block.Data, timebytes, block.Hash}, []byte{})
	// 4.生成hash
	HashValue := sha256.Sum256(blockbytes)

	// 处理 HashValue是32字节
	block.Hash = HashValue[:]
} */

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
	//fmt.Println("old block = ", block)
	// 设置hash
	//block.SetHash()

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
func CreateGenesisBlock(data string) *Block {
	// 高度 hash可知
	var csqk [32]byte
	return NewBlock(data, 1, csqk[:])
}
