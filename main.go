package main

import (
	"fmt"

	block "iniyou.com/BLOCK"
)

func main() {
	var csqk [16]byte
	// block := block.NewBlock("Genenis Block", 1, []byte{0})
	block := block.NewBlock("Genenis Block", 1, csqk[:])
	fmt.Println("newblock = ", block)
}
