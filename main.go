package main

import (
	"fmt"

	block "iniyou.com/BLOCK"
)

func main() {
	// genesisblock
	//genesisBlockchain := block.CreateBlockchainWithGenesisBlock()

	/* 	// newblock
	   	genesisBlockchain.AddBlockToBlockchain("1=>2", genesisBlockchain.Blocks[len(genesisBlockchain.Blocks)-1].Height+1, genesisBlockchain.Blocks[len(genesisBlockchain.Blocks)-1].Hash)

	   	genesisBlockchain.AddBlockToBlockchain("2=>3", genesisBlockchain.Blocks[len(genesisBlockchain.Blocks)-1].Height+1, genesisBlockchain.Blocks[len(genesisBlockchain.Blocks)-1].Hash)

	   	genesisBlockchain.AddBlockToBlockchain("3=>4", genesisBlockchain.Blocks[len(genesisBlockchain.Blocks)-1].Height+1, genesisBlockchain.Blocks[len(genesisBlockchain.Blocks)-1].Hash)
	   	fmt.Println(genesisBlockchain) // &{[0xc000076060]}
	   	fmt.Println("block0=>", genesisBlockchain.Blocks[0])
	   	fmt.Println("block2=>", genesisBlockchain.Blocks[1])
	   	fmt.Println("block3=>", genesisBlockchain.Blocks[2])
	   	fmt.Println("block4=>", genesisBlockchain.Blocks[3]) */

	block1 := block.NewBlock("text", 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	fmt.Println("block.nonce=", block1.Nonce)
	fmt.Println("block.hash =", block1.Hash)

	// 上面区块已经验证完了，快速验证
	pow := block.NewProofOfWork(block1)
	fmt.Println("is vaied", pow.IsVaild())
}
