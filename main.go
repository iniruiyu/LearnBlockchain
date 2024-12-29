package main

import (
	"fmt"

	block "iniyou.com/BLOCK"
)

func main() {
	genesisBlockchain := block.CreateBlockchainWithGenesisBlock()
	fmt.Println(genesisBlockchain) // &{[0xc000076060]}
	fmt.Println(genesisBlockchain.Blocks[0])
	// {1 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	// [71 101 110 101 115 105 115 32 98 108 111 99 107 32 68 97 116 97 46 46 46]
	// 1723208163
	// [243 105 111 102 243 102 209 248 15 69 240 134 180 66 147 40 80 118 59 31 68 147 41 56 105 152 50 26 138 189 23 254]}
}
