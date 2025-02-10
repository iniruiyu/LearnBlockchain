package cli

import (
	"fmt"
	"os"

	block "iniyou.com/blockchain"
)

func (cli *Cli) send(from []string, to []string, amount []string) {
	if !block.DbExists() {
		fmt.Println("数据库不存在")
		os.Exit(1)
		return
	}
	blockchain := block.GetBlockObject()
	defer blockchain.DB.Close()
	blockchain.MineNewBlock(from, to, amount)
}
