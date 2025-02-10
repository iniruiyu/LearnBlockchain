package cli

import (
	"fmt"
	"os"

	block "iniyou.com/blockchain"
)

// 查询地址余额
func (cli *Cli) getBalance(address string) {
	if !block.DbExists() {
		fmt.Println("数据库不存在")
		os.Exit(1)
		return
	}
	blockchain := block.GetBlockObject()
	defer blockchain.DB.Close()
	fmt.Println(address, "Address Amount", blockchain.GetBalance(address))
}
