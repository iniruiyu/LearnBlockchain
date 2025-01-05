package main

import (
	block "iniyou.com/BLOCK"
)

// persistance
// 数据持久化

func main() {
	// genesisblock
	//Blockchain := block.CreateBlockchainWithGenesisBlock()
	//defer Blockchain.DB.Close()

	Cli := block.Cli{}
	Cli.Run()
}
