package main

import (
	block "iniyou.com/BLOCK"
)

// persistance
// 数据持久化

func main() {
	// genesisblock
	genesisBlockchain := block.CreateBlockchainWithGenesisBlock()
	defer genesisBlockchain.DB.Close()
	// newblock
	genesisBlockchain.AddBlockToBlockchain("1=>2")

	genesisBlockchain.AddBlockToBlockchain("2=>3")
	genesisBlockchain.AddBlockToBlockchain("3=>4")
	genesisBlockchain.AddBlockToBlockchain("4=>5")
	//genesisBlockchain.AddBlockToBlockchain("3=>4", genesisBlockchain.Blocks[len(genesisBlockchain.Blocks)-1].Height+1, genesisBlockchain.Blocks[len(genesisBlockchain.Blocks)-1].Hash)
	/*fmt.Println(genesisBlockchain) // &{[0xc000076060]}
	 fmt.Println("block0=>", genesisBlockchain.Blocks[0])
	fmt.Println("block2=>", genesisBlockchain.Blocks[1])
	fmt.Println("block3=>", genesisBlockchain.Blocks[2])
	fmt.Println("block4=>", genesisBlockchain.Blocks[3]) */
	genesisBlockchain.PrintBlockchain()
}
