package block

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Cli struct{}

func (cli *Cli) Run() {
	isVaildArgs()
	printchainCMD := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockCMD := flag.NewFlagSet("addblock", flag.ExitOnError)
	flagAddBlockData := addBlockCMD.String("data", "iniyou.com", "交易数据....")

	createBlockChainCMD := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	createBlockChainData := createBlockChainCMD.String("data", "Genesis Block daata...", "创世区块交易数据....")
	switch os.Args[1] {
	case "addblock":
		err := addBlockCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		cli.AddBlock(*flagAddBlockData)
	case "printchain":
		err := printchainCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		cli.Printchain()
	case "createBlockchain":
		err := createBlockChainCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		//
		cli.CreateGenesisBlockchain(*createBlockChainData)
	default:
		printUsage()
		os.Exit(1)
	}
	if createBlockChainCMD.Parsed() {
		if *createBlockChainData == "" {
			fmt.Println("交易数据不能为空")
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*createBlockChainData)
	}
	if addBlockCMD.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*flagAddBlockData)
	}
	if printchainCMD.Parsed() {
		fmt.Println(*printchainCMD, "交易数据")
	}

}
func (cli *Cli) AddBlock(data string) {
	if !DbExists() {
		fmt.Println("数据库不存在")
		os.Exit(1)
		return
	}
	blockchain := GetBlockObject()
	defer blockchain.DB.Close()
	blockchain.AddBlockToBlockchain(data)
}
func (cli *Cli) Printchain() {
	if !DbExists() {
		fmt.Println("数据库不存在")
		os.Exit(1)
		return
	}
	blockchain := GetBlockObject()
	defer blockchain.DB.Close()
	blockchain.PrintBlockchain()
}
func (cli *Cli) CreateGenesisBlockchain(data string) {
	fmt.Println(data)
	CreateBlockchainWithGenesisBlock(data)

	blockchain := GetBlockObject()
	defer blockchain.DB.Close()
	blockchain.PrintBlockchain()
}
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\t createBlockcahin -data DATA -交易数据")
	fmt.Println("\t addblock -data DATA -交易数据")
	fmt.Println("\t printchain -输出区块信息")
}
func isVaildArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
