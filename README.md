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
  

### 6.1 code

* 
