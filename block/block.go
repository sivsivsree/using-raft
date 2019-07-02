package block

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type Block struct {
	index     int
	timestamp string
	data      string
	hash      string
	prevHash  string
}

var Blockchain []Block

var mutex = &sync.Mutex{}

func AddNewBlock(data string) ([]Block, error) {

	var newBlock Block

	t := time.Now()
	mutex.Lock()
	oldBlock := Blockchain[len(Blockchain)-1]
	newBlock.index = oldBlock.index + 1
	newBlock.timestamp = t.String()
	newBlock.data = data
	newBlock.prevHash = oldBlock.hash
	newBlock.hash = calculateHash(newBlock)
	Blockchain = append(Blockchain, newBlock)
	mutex.Unlock()
	return Blockchain, nil

}

func InitGenesis() ([]Block, error) {
	if len(Blockchain) < 1 {
		genesisBlock := Block{0, time.Now().String(), "", calculateHash(Block{}), "0"}
		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}
	return Blockchain, nil
}

func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.index+1 != newBlock.index {
		return false
	}

	if oldBlock.hash != newBlock.prevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.hash {
		return false
	}

	return true
}

func calculateHash(block Block) string {
	record := string(block.index) + block.timestamp + string(block.data) + block.prevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
