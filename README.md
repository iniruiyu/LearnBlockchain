# LearnBlockchain

## 1.Create block struct

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

## 2. Genesis Block

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

