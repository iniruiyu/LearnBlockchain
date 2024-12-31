package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 0000 0000 0000 0000 1001 0001 0000 .... 0001
// 256 位Hash 里面，前面至少要有16个0
const targetBit = 16

type ProofOfWork struct {
	Block  *Block   // 要去验证的区块
	Target *big.Int // 大数据存储
	// int64可能会溢出 他范围比较大
	// 代表我们数据的难度
}

// 数据拼接，返回字节数组
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(pow.Block.Height)),
			IntToHex(int64(targetBit)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// 判断hash是否有效
func (proofOfWork *ProofOfWork) IsVaild() bool {
	/* 	proofOfWork.Block.Hash
	   	proofOfWork.Target */
	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.Block.Hash)
	/*Cfunc (x *big.Int) Cmp(y *big.Int) (r int)
	mp compares x and y and returns:
		-1 if x <  y
		 0 if x == y
		+1 if x >  y */
	return proofOfWork.Target.Cmp(&hashInt) == 1
}

func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {
	// 1. 将Block的属性，拼接成字节数组
	//dataBytes := proofOfWork.prepareData()
	// 2. 生成hash

	// 3. 判断hash有效性，如果满足条件，跳出循环

	var hashInt big.Int // 存储我们新生成的hash值
	var hash [32]byte
	for nonce := 0; ; nonce++ {
		// 1. 准备数据
		dataBytes := proofOfWork.prepareData(nonce)
		// 2. 生成hash
		hash = sha256.Sum256(dataBytes)
		//fmt.Printf("\r%x\n", hash)
		fmt.Printf("\r%x", hash)
		// 2.2. 存储到HashInt
		hashInt.SetBytes(hash[:])
		// 3. 判断hash有效性，如果满足条件，跳出循环

		/*
			func (x *big.Int) Cmp(y *big.Int) (r int)
			Cmp compares x and y and returns:

				-1 if x <  y
				 0 if x == y
				+1 if x >  y
		*/
		if proofOfWork.Target.Cmp(&hashInt) == 1 {
			fmt.Println()
			return hash[:], int64(nonce)
		}
	}
	//return nil, 0
}

// 创建新的工作量证明对象
func NewProofOfWork(Block *Block) *ProofOfWork {
	// 1. big.Int对象
	// 1 255个0 最后一个为1
	// 保证前面有16个零，可以进行左移位
	/* target 2个0
	0000 0001
	左移 8-2 =6 位
	0100 0000  =64
	0010 0000  =32
	<32就符合
	*/

	// 1.创建一个初始值为1的taget
	target := big.NewInt(1)
	// 2.左移256 - targetBit
	target = target.Lsh(target, 256-targetBit)

	return &ProofOfWork{Block: Block, Target: target}
}
