package cli

import block "iniyou.com/blockchain"

// 创建创世区块
func (cli *Cli) CreateGenesisBlockchain(address string) {

	// 1.创建
	blockchain := block.CreateBlockchainWithGenesisBlock(address)
	defer blockchain.DB.Close()

	/*blockchain := GetBlockObject()
	defer blockchain.DB.Close()*/
	//blockchain.PrintBlockchain()
}
