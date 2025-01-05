# LearnBlockchain

## 1. Create block struct

* block struct

  * ```go
    type Block struct {
    	// 1. block height
    	Height int64
    	// 2. Previous block hash
    	PrevBlockHash []byte
    	// 3. Data,In the future, it will be used to store data such as transactions
    	Data []byte
    	// 4. Timestamp
    	Timestamp int64
    	// 5. Block Hash
    	Hash []byte
    }
    ```

* Set Block Hash Function

  * ```go
    // 2. Set Block hash
    func (block *Block) SetHash() {
    
    	// 1. Convert Height to a byte array
    	// use IntToHex() Function
    	heightBytes := IntToHex(block.Height)
    	// 2.Converts Timestamp to a byte array
    	// 2.1  strconv.FormatInt()  
    	// The first parameter converts int64 to a string
    	// The second parameter ranges from 2 to 36, representing the base system
    	timeString := strconv.FormatInt(block.Timestamp, 2)
    	timebytes := []byte(timeString)
    	fmt.Println("SetHash timeString", timeString, "\n timebytes := ", timebytes, "\n heightBytes:=", heightBytes)
    	// 3.Concatenate all properties
    	blockbytes := bytes.Join([][]byte{heightBytes, block.PrevBlockHash, block.Data, timebytes, block.Hash}, []byte{})
    	// 4.Generating hash
    	HashValue := sha256.Sum256(blockbytes)
    
    	// Processing HashValue is 32 bytes
    	block.Hash = HashValue[:]
    }
    
    // utils.go
    // Convert int64 to byte Array
    func IntToHex(num int64) []byte {
    	buff := new(bytes.Buffer)
    	err := binary.Write(buff, binary.BigEndian, num)
    	if err != nil {
    		log.Panic("IntToHex error", err)
    	}
    	return buff.Bytes()
    }
    ```

  * 
* Create block function

  * ```go
    func NewBlock(data string, height int64, prevBlockHash []byte) *Block {
    	// CreateBlock
    	block := &Block{
    		Height:        height,
    		PrevBlockHash: prevBlockHash,
    		Data:          []byte(data),
    		Timestamp:     time.Now().Unix(),
    		Hash:          nil,
    	}
    	fmt.Println("old block = ", block)
    	// Sethash
    	block.SetHash()
    	return block
    }
    ```

## 2. Create Genesis Block

* The Genesis block is the first block of the blockchain. 

* It usually has a **height of 1** and a **block Hash** of **0 as a 32-bit array**

  * ```go
    // Create Genesis Block
    func CreateGenesisBlock(data string) *Block {
    	// create [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]
    	// 32-bit byte array 
    	var ZeroBlockHash [32]byte
    	return NewBlock(data, 1, ZeroBlockHash[:])
    }
    
    // main.go
    block := block.CreateGenesisBlock("Genesis Block")
    fmt.Println("Genesis block = ", block)

## 3. Add Genesis Block to Blockchain

* A blockchain is composed of many blocks, and the genesis block is the first block of the blockchain.

* Blockchain.go

  * Blockchain struct

    * ```go
    	type Blockchain struct {
    		Blocks []*Block // Stores ordered blocks
    	}
    	```

  * Create blockchain with Genesis block

    * ```go
      func CreateBlockchainWithGenesisBlock() *Blockchain {
      	GenesisBlock := CreateGenesisBlock("Genesis block Data...")
      	return &Blockchain{Blocks: []*Block{GenesisBlock}}
      }
      ```

* Add block to blckchain

  * ```go
    func (blockchain *Blockchain) AddBlockToBlockchain(data string, height int64, prevHash []byte) {
    	newBlock := NewBlock(data, height, prevHash)
    	// Adds blocks to the chain array
    	blockchain.Blocks = append(blockchain.Blocks, newBlock)
    }
    ```

    

## 4. Proof Of Work(POW)

* **Proof of Work (PoW)** is the earliest consensus algorithm adopted in blockchain technology.

* How It Works

  * <details><summary>Puzzle Solving</summary>难题计算</details>

  * <details><summary>Simple Verification</summary>验证简单</details>

  * <details><summary>Reward Mechanism</summary>奖励机制</details>



* Create Proof Of Work Struct

  * ```go
    type ProofOfWork struct {
            Block *Block // The block to validate
            Target *big.Int // Big Data Storage
            // Represents the difficulty of our data
            // int64 may overflow its range larger
    }
    ```

* Difficulty

  * After running, 0000 0000 0000 0000 1001 0001 0000 .... 0001 is shifted **256-targetBit** to the left

  * ```go
    // Difficulty
    // 0000 0000 0000 0000 1001 0001 0000 .... 0001
    // A 256-bits Hash must have at least 16/targetBit zeroes in front of it
    const targetBit = 16
    ```

* Run()

  1. Concatenate the properties of the Block into a byte array
  2. Generate hash
  3. If the hash is valid, the result is returned

  * ```go
    func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {
    	var hashInt big.Int // Store our newly generated hash value
    	var hash [32]byte
    	for nonce := 0; ; nonce++ {
    		// 1. prepare Data
    		dataBytes := proofOfWork.prepareData(nonce)
    		// 2. Generate hash
    		hash = sha256.Sum256(dataBytes)
    		fmt.Printf("\r%x", hash)
    		// 2.2. Store in HashInt
    		hashInt.SetBytes(hash[:])
    		// 3. Checking the Validity of generate hash
    
    		/*	func (x *big.Int) Cmp(y *big.Int) (r int)
    			Cmp compares x and y and returns:
    				-1 if x <  y
    				 0 if x == y
    				+1 if x >  y
    		*/
    		if proofOfWork.Target.Cmp(&hashInt) == 1 {
    			fmt.Println()
    			return hash[:], int64(nonce)
    		}
    	}
    }
    
    // Concatenate the properties of the Block into a byte array.
    func (pow *ProofOfWork) prepareData(nonce int) []byte {
    	data := bytes.Join(
    		[][]byte{
    			pow.Block.PrevBlockHash,
    			pow.Block.Data,
    			IntToHex(pow.Block.Timestamp),
    			IntToHex(int64(pow.Block.Height)),
    			IntToHex(int64(targetBit)),
    			IntToHex(int64(nonce)),
    		},
    		[]byte{},
    	)
    	return data
    }
    ```

  * Create a new proof of work object

    ```go
    func NewProofOfWork(Block *Block) *ProofOfWork {
    	/* target two 0
    	0000 0001
    	Shift left(8-2 =6) bit
    	0100 0000  =64
    	0010 0000  =32
    	< Just move it two places to the left ,32
    	*/
    	// 1.Create a taget with an initial value of 1
    	target := big.NewInt(1)
    	// 2.Shift 256-target Bit to the left
    	target = target.Lsh(target, 256-targetBit)
    	return &ProofOfWork{Block: Block, Target: target}
    }
    ```

* Verify that the hash is valid
  ```go
  // Determine whether the generated hash is preceded by Target zeros
  func (proofOfWork *ProofOfWork) IsVaild() bool {
  	/* 	proofOfWork.Block.Hash
  	   	proofOfWork.Target */
  	var hashInt big.Int
  	hashInt.SetBytes(proofOfWork.Block.Hash)
  	/* func (x *big.Int) Cmp(y *big.Int) (r int)
  	mp compares x and y and returns:
  		-1 if x <  y
  		 0 if x == y
  		+1 if x >  y */
  	return proofOfWork.Target.Cmp(&hashInt) == 1
  }
  
  func main() {
  	block1 := block.NewBlock("text", 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
  	fmt.Println("block.nonce=", block1.Nonce)
  	fmt.Println("block.hash =", block1.Hash)
  	// The upper block has been verified, quick verification
  	pow := block.NewProofOfWork(block1)
  	fmt.Println("is vaied", pow.IsVaild())	//true
  }
  ```

  

## 5. Serialize

* To store blocks on disk, serialization is needed

* block.go

* ```go
  // use "encoding/gob" 
  // Serialize the block into a byte array
  func (block *Block) Serialize() []byte {
  	var result bytes.Buffer
  	encoder := gob.NewEncoder(&result)
  	err := encoder.Encode(block)
  	if err != nil {
  		log.Panic(err)
  	}
  	return result.Bytes()
  }
  
  // Deserializing the byte array returns the block structure
  func DeSerializeBlock(blockBytes []byte) *Block {
  	var block Block
  	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
  	err := decoder.Decode(&block)
  	if err != nil {
  		log.Panic(err)
  	}
  	return &block
  }
  ```



## 6. DataBase

* "github.com/boltdb/bolt"

* Introduction

  * CreateTable
  
  * ```go
    package main
    
    import (	"github.com/boltdb/bolt")
    
    func main() {
    	// Open the my.db data file in your current directory.
    	// It will be created if it doesn't exist.
    	db, err := bolt.Open("my.db", 0600, nil)
    	if err != nil {
    		log.Fatal(err)
    	}
    	err = db.Update(func(tx *bolt.Tx) error {
    		// 1. create BlockBucket Table
    		b, err := tx.CreateBucket([]byte("BlockBucket"))
            
    		if err != nil {return fmt.Errorf("create bucket error,%s", err)}
    		// 2. Putting data into a table, in key-value pairs
    		if b != nil {
                err := b.Put([]byte("1"), []byte("Send 100 to 2"))
    			if err != nil {log.Panic("Save data to bucket error", err)}
    		}
    		return nil
    	})
    	defer db.Close()
    
    }
    ```
  
  * GetTable
  
    * ```go
      func main() {
      	// Open the my.db data file in your current directory.
      	// It will be created if it doesn't exist.
      	db, err := bolt.Open("my.db", 0600, nil)
      	if err != nil {
      		log.Fatal(err)
      	}
      	err = db.Update(func(tx *bolt.Tx) error {
      		// 1. Get BlockBucket Table
      		b := tx.Bucket([]byte("BlockBucket"))
      
      		// 2. Putting data into a table, in key-value pairs
      		if b != nil {
      			err := b.Put([]byte("1"), []byte("Send 100 to 2"))
      			if err != nil {log.Panic("Save data to bucket error", err)}
      		}
      		return nil
      	})
      	defer db.Close()
      }
      ```
  
  * ViewTable
  
    * ```go
      	// ViewTable
      	err = db.View(func(tx *bolt.Tx) error {
      		b := tx.Bucket([]byte("BlockBucket"))
      		if b != nil {
      			key1 := b.Get([]byte("1"))
      			fmt.Println(key1)
      		}
      		return nil
      	})
      ```
  

## 7. Save Genesis Block To BoltDB

* Persistence store

* Blockchain uses db object, just need to add db object to blockchain property.

* Modify Blockchain Struct

  * ```go
    const dbName = "blockchain.db"
    const blockTableName = "block"
    
    type Blockchain struct {
    	Tip []byte   // The hash of the latest block
    	DB  *bolt.DB // Database
    }
    ```

* Modify CreateBlockchainWithGenesisBlock() Function

  * ```go
    func CreateBlockchainWithGenesisBlock() *Blockchain {
    	var blockHash []byte
    	// Create or open the database
    	db, err := bolt.Open(dbName, 0600, nil)
    	if err != nil {
    		log.Fatal(err)
    	}
    	defer db.Close()
    	err = db.Update(func(tx *bolt.Tx) error {
    		b, err := tx.CreateBucket([]byte(blockTableName))
    		if err != nil {
    			log.Panic(err)
    		}
    		//The table exists
    		if b != nil {
    			// Create GenisisBlock
    			GenesisBlock := CreateGenesisBlock("Genesis block Data...")
    			// Store the Genesis block in a table
    			err := b.Put(GenesisBlock.Hash, GenesisBlock.Serialize())
    			if err != nil {
    				log.Panic(err)
    			}
    			// Store the hash of the latest block
    			err = b.Put([]byte("TipBlockHash"), GenesisBlock.Hash)
    			if err != nil {
    				log.Panic(err)
    			}
    			blockHash = GenesisBlock.Hash
    		}
    		return nil
    	})
    	return &Blockchain{blockHash, db}
    }
    ```

## 8. Add New Block To BoltDb

* ```go
  func (blockchain *Blockchain) AddBlockToBlockchain(data string) {
  	blockchain.DB.Update(func(tx *bolt.Tx) error {
  		// 1. Get Bucket
  		b := tx.Bucket([]byte(blockTableName))
  		// 2. Create New Block
  		if b != nil {
  			// 1. **Get The Latest Block**
  			TipBlockBytes := b.Get(blockchain.Tip)
  			// 2. DeSerialize
  			TipBlock := DeSerializeBlock(TipBlockBytes)
  			// 3. Create New Block 
  			newBlock := NewBlock(data, TipBlock.Height+1, TipBlock.Hash)
  			// 4. Serialize New Block ,place it in the obtained table
  			err := b.Put(newBlock.Hash, newBlock.Serialize())
  			if err != nil {
  				log.Panic(err)
  			}
  			// Store the hash of the latest block
  			err = b.Put([]byte("TipBlockHash"), newBlock.Hash)
  			if err != nil {
  				log.Panic(err)
  			}
  		}
  		return nil
  	})
  
  }
  ```

  

## 9. Print Blockchain Info

* Iterate over all blocks to output information
  ```go
  func (blockchain *Blockchain) PrintBlockchain() {
  
  	var block *Block
  	var currentHash []byte = blockchain.Tip
  	for {
  		err := blockchain.DB.Update(func(tx *bolt.Tx) error {
  			// 1.Get Bucket
  			b := tx.Bucket([]byte(blockTableName))
  			if b != nil {
  				// Gets the byte array of the current block
  				blockBytes := b.Get(currentHash)
  				// DeSerializeBlock
  				block = DeSerializeBlock(blockBytes)
  				fmt.Printf("Height:%d\n", block.Height)
  				fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
  				fmt.Printf("Data:%s\n", block.Data)
  				//fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
  				fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 PM"))
  				fmt.Printf("Hash:%x\n", block.Hash)
  				fmt.Printf("Nonce:%d\n\n", block.Nonce)
  			}
  			return nil
  		})
  		if err != nil {
  			log.Panic(err)
  		}
  		var hashInt big.Int
  		hashInt.SetBytes(block.PrevBlockHash)
  		if big.NewInt(0).Cmp(&hashInt) == 0 {
  			break
  		}
  		// Iterate
  		currentHash = block.PrevBlockHash
  	}
  }
  ```



## 10. Iterator

Write an iterator that optimizes the above code to reduce repetition

* Iterator Struct

* ```go
  type BlockchainIterator struct {
  	CurrentHash []byte   // The latest Block Hash
  	DB          *bolt.DB // DB
  }
  // Get the next block
  func (blockchainIterator *BlockchainIterator) NextPrevBlock() *Block {
  	var currentBlock *Block
  	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
  		b := tx.Bucket([]byte(blockTableName))
  		if b != nil {
  			currentBlockBytes := b.Get(blockchainIterator.CurrentHash)
  			// The block corresponding to the current iterator Hash is retrieved
  			currentBlock = DeSerializeBlock(currentBlockBytes)
  			blockchainIterator.CurrentHash = currentBlock.PrevBlockHash
  		}
  		return nil
  	})
  	if err != nil {
  		log.Panic(err)
  	}
  	return currentBlock
  }
  // Get the blockchain iterator object
  func (blockchain *Blockchain) Iterator() *BlockchainIterator {
  	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
  }
  ```

* Modify PrintBlockchain()

  * ```go
    func (blockchain *Blockchain) PrintBlockchain() {
    
    	blockchainIterator := blockchain.Iterator()
    	for {
    		block := blockchainIterator.NextPrevBlock()
    		fmt.Printf("Height:%d\n", block.Height)
    		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
    		fmt.Printf("Data:%s\n", block.Data)
    		fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 PM"))
    		fmt.Printf("Hash:%x\n", block.Hash)
    		fmt.Printf("Nonce:%d\n\n", block.Nonce)
    		var hashInt big.Int
    		hashInt.SetBytes(block.PrevBlockHash)
    		if big.NewInt(0).Cmp(&hashInt) == 0 {
    			break
    		}
    	}
    }
    ```



## 11. Cli

### 11.1 os.Args

* Input
  * main.exe addBlock -data "block1"

* ```go
  // main.exe addBlock -data "block1"
  func main() {
  	args := os.Args
  	fmt.Printf("%v\n", args)		// [D:\dev\main.exe addblock -data block1]
  	fmt.Printf("%v\n", args[1])		// addblock
  	fmt.Printf("%v\n", args[2])		// -data
  	fmt.Printf("%v\n", args[2:])	// [-data block1]
  }
  ```



### 11.2 flag

  * input

  * ```go
    flagString := flag.String("printchain", "string", "输出所有的区块信息")
    flagInt := flag.Int("number", 5, "输出一个整数")
    flagBool := flag.Bool("open", false, "判断真假")
    flag.Parse()
    fmt.Printf("%s\n", *flagString)		// string
    fmt.Printf("%d\n", *flagInt)		// 5
    fmt.Printf("%v\n", *flagBool)		// false
    ```



### 11.3 Cli Struct

Cli Struct

```go
type Cli struct{}
```

Cli function

* Check OsArgs Valid

  * ```go
    func isVaildArgs() {
    	if len(os.Args) < 2 {
    		printUsage()
    		os.Exit(1)
    	}
    }
    ```

* Print Cli Cmd

  * ```go
    func printUsage() {
    	fmt.Println("Usage:")
    	fmt.Println("\t createBlockchain -address -地址")
    	fmt.Println("\t send -from FROM -to TO -amount AMOUNT -交易")
    	fmt.Println("\t printchain -输出区块信息")
    	fmt.Println("\t getbalance -address -获取余额")
    }
    ```

* ```go
  func (cli *Cli) Run() {
  	isVaildArgs()
  	// 1.printchain
  	printchainCMD := flag.NewFlagSet("printchain", flag.ExitOnError)
  	// 2.send
  	//.\main.exe send -from '[\"address1","address2\"]' -to '[\"address3\",\"address4\"]' -amount '[\"2\",\"3\"]'
  	sendBlockCMD := flag.NewFlagSet("send", flag.ExitOnError)
  	flagFrom := sendBlockCMD.String("from", "", "from address....")
  	flagTo := sendBlockCMD.String("to", "", "to address....")
  	flagAmount := sendBlockCMD.String("amount", "", "amount....")
  
  	//3.createblockchain
  	createBlockChainCMD := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
  	createBlockChainWithAddress := createBlockChainCMD.String("address", "", "Genesis Block Address....")
  
  	//4.getbalance
  	getbalanceCMD := flag.NewFlagSet("getbalance", flag.ExitOnError)
  	getbalanceCMDWithAddress := getbalanceCMD.String("address", "", "Get Address Amount")
  
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
  			fmt.Println("Address == "")
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
  		from := block.JSONToArray(*flagFrom)
  		to := block.JSONToArray(*flagTo)
  		amount := block.JSONToArray(*flagAmount)
  		cli.send(from, to, amount)
  	}
  	if printchainCMD.Parsed() {
  		cli.Printchain()
  	}
  	if getbalanceCMD.Parsed() {
  		if *getbalanceCMDWithAddress == "" {
  			fmt.Println("Address == "")
  			printUsage()
  			os.Exit(1)
  		}
  		cli.getBalance(*getbalanceCMDWithAddress)
  	}
  
  }
  ```



### 11.4 Cli Function

1. Create Genesis Blockchain

   * ```go
     func (cli *Cli) CreateGenesisBlockchain(address string) {
     	blockchain := block.CreateBlockchainWithGenesisBlock(address)
     	defer blockchain.DB.Close()
     }
     ```

   * ```go
     // creates a blockchain with a genesis block
     func CreateBlockchainWithGenesisBlock(genesisBlockAddress string) *Blockchain {
     	// Exit if the database exists
     	if DbExists() {
     		fmt.Println("Genesis block already exists")
     		os.Exit(1)
     		return nil
     	}
     	// Create or open the database
     	db, err := bolt.Open(dbName, 0600, nil)
     	if err != nil {
     		log.Fatal(err)
     	}
     	// Close the database when the function exits
     	defer db.Close()
     
     	var GenesisBlockHash []byte
     	err = db.Update(func(tx *bolt.Tx) error {
     
     		// Create a bucket
     		b, err := tx.CreateBucket([]byte(blockTableName))
     		if err != nil {
     			log.Panic(err)
     		}
     		// Bucket already exists
     		if b != nil {
     			// Create the GenesisBlock
     			// Create a Coinbase Transaction
     			txCoinbase := NewCoinbaseTransaction(genesisBlockAddress)
     			GenesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
     			// Store the genesis block in the bucket
     			err := b.Put(GenesisBlock.Hash, GenesisBlock.Serialize())
     			if err != nil {
     				log.Panic(err)
     			}
     			// Store the hash of the latest block
     			err = b.Put([]byte("TipBlockHash"), GenesisBlock.Hash)
     			if err != nil {
     				log.Panic(err)
     			}
                 // Use to Return
     			GenesisBlockHash = GenesisBlock.Hash
     		}
     
     		return nil
     	})
     	if err != nil {
     		log.Panic(err)
     	}
     	return &Blockchain{GenesisBlockHash, db}
     }
     ```

2. Print Block Info

   * ```
     func (cli *Cli) Printchain() {
     	if !block.DbExists() {
     		fmt.Println("Database does not exist")
     		os.Exit(1)
     		return
     	}
     	blockchain := block.GetBlockObject()
     	defer blockchain.DB.Close()
     	blockchain.PrintBlockchain()
     }
     ```

   * ```go
     // PrintBlockchain prints information for all blocks in the blockchain
     func (blockchain *Blockchain) PrintBlockchain() {
     	// Initialize the iterator to start from the latest block
     	blockchainIterator := blockchain.Iterator()
     	// Iterate over all blocks in the blockchain
     	for {
     		// Get the next previous block in the chain
     		block := blockchainIterator.NextPrevBlock()
     		// Print the block header information
     		fmt.Printf("---------------Block Height %d------------------\n", block.Height)
     		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
     		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 PM"))
     		fmt.Printf("Hash: %x\n", block.Hash)
     		fmt.Printf("Nonce: %d\n", block.Nonce)
     
     		// Print the transactions in the block
     		fmt.Printf("---------------Block %d - Txs: %v---------------\n", block.Height, len(block.Txs))
     		for _, tx := range block.Txs {
     			fmt.Printf("Txs: tx.hash=%x\n", tx.TxHash)
     			fmt.Printf("Vins →   ")
     			for _, in := range tx.Vins {
     				fmt.Printf("in.TxHash: %x, Signature: %x, (in.Vout Index): %d\n",
     					in.TxHash, in.ScriptSig, in.Vout)
     			}
     			fmt.Printf("Vouts→   \n")
     			for i, out := range tx.Vouts {
     				fmt.Printf("tx.%d Value: %d, ScriptPublicKey: %s\n", i, out.Value, out.ScriptPublicKey)
     			}
     			fmt.Printf("--------------End of Block %d - Tx Ended--------------\n", block.Height)
     		}
     		// If the current block's previous hash is zero, it means we've reached the genesis block
     		if bytes.Equal(block.PrevBlockHash, []byte{}) {
     			break
     		}
     	}
     }
     ```

###  That's pretty easy, now we're going to write some trading code.



## 12. Transaction

