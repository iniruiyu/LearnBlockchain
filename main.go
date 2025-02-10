package main

import (
	cli "iniyou.com/cli"
)

// persistance
// 数据持久化

func main() {
	// genesisblock
	//Blockchain := block.CreateBlockchainWithGenesisBlock()
	//defer Blockchain.DB.Close()

	Cli := cli.Cli{}
	Cli.Run()
}
