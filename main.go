package main

import (
	"fmt"

	block "iniyou.com/BLOCK"
)

func main() {
	//var csqk [32]byte
	// block := block.NewBlock("Genenis Block", 1, []byte{0})
	//block := block.NewBlock("Genenis Block", 1, csqk[:])

	block := block.CreateGenesisBlock("Genesis Block")
	fmt.Println("newblock = ", block)
}
