package blockchain

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
)

// 将int64转化为字节数组
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic("IntToHex error", err)
	}
	return buff.Bytes()
}

// 标准的JSON 字符串转数组
func JSONToArray(jsonString string) []string {
	//json到string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		fmt.Println(" jsonString=", jsonString, "err= ", err)
		log.Panic(err)
	}
	return sArr

}
