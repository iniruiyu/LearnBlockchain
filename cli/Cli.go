package cli

import (
	"flag"
	"fmt"
	"log"
	"os"

	block "iniyou.com/blockchain"
)

type Cli struct{}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\t createBlockchain -address -地址")
	fmt.Println("\t send -from FROM -to TO -amount AMOUNT -交易")
	fmt.Println("\t printchain -输出区块信息")
	fmt.Println("\t getbalance -address -获取余额")
}
func isVaildArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *Cli) Run() {
	isVaildArgs()
	// 1.printchain
	printchainCMD := flag.NewFlagSet("printchain", flag.ExitOnError)

	// 2.send
	//.\main.exe send -from '[\"liyuechun","zhangqiang\"]' -to '[\"junchen\",\"xiaoyong\"]' -amount '[\"2\",\"3\"]'
	sendBlockCMD := flag.NewFlagSet("send", flag.ExitOnError)
	flagFrom := sendBlockCMD.String("from", "", "转账源地址....")
	flagTo := sendBlockCMD.String("to", "", "转账目的地址....")
	flagAmount := sendBlockCMD.String("amount", "", "转账金额....")

	//3.createblockchain
	createBlockChainCMD := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	createBlockChainWithAddress := createBlockChainCMD.String("address", "", "创建创世区块的地址....")

	//4.getbalance
	getbalanceCMD := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getbalanceCMDWithAddress := getbalanceCMD.String("address", "", "获取余额的地址")

	switch os.Args[1] {
	case "send":
		err := sendBlockCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printchainCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createBlockchain":
		err := createBlockChainCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getbalanceCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
	if createBlockChainCMD.Parsed() {
		if *createBlockChainWithAddress == "" {
			fmt.Println("Address不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.CreateGenesisBlockchain(*createBlockChainWithAddress)
	}
	if sendBlockCMD.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}
		//fmt.Println(*flagFrom)

		//fmt.Println(block.JSONToArray(*flagFrom))
		//fmt.Println(block.JSONToArray(*flagTo))
		//fmt.Println(block.JSONToArray(*flagAmount))
		//cli.AddBlock(*flagAddBlockData)
		from := block.JSONToArray(*flagFrom)
		to := block.JSONToArray(*flagTo)
		amount := block.JSONToArray(*flagAmount)
		cli.send(from, to, amount)
	}
	if printchainCMD.Parsed() {
		//fmt.Println(*printchainCMD, "交易数据")
		cli.Printchain()
	}
	if getbalanceCMD.Parsed() {
		if *getbalanceCMDWithAddress == "" {
			fmt.Println("地址不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceCMDWithAddress)
	}

}
