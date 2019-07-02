package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/sivsivsree/using-raft/logstore"
	bolt "go.etcd.io/bbolt"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Pos           int
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func NewBlock(data string, prevBlock *Block) *Block {
	pos := prevBlock.Pos + 1
	block := &Block{pos, time.Now().Unix(), []byte(data), prevBlock.Hash, []byte{}}
	block.SetHash()
	return block
}

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	// AddBlockToTheBucket
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock)

	newBlockBytes, _ := json.Marshal(newBlock)

	err := logstore.AddNewBlockToBucket(newBlock.Pos, newBlockBytes)

	if err != nil {
		log.Fatal(err)
	}

	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGenesisBlock() *Block {
	// Add NewGenesis to the Block
	pos := 0
	genisis := &Block{pos, time.Now().Unix(), []byte("Geneisis"), []byte(""), []byte{}}
	genisis.SetHash()
	newBlockBytes, _ := json.Marshal(genisis)
	err := logstore.AddNewBlockToBucket(0, newBlockBytes)
	if err != nil {
		log.Fatal(err)
	}
	return genisis
}

func NewBlockchain() *Blockchain {
	// createBucketIfNotExists
	err := logstore.CreateBucketIfNotExist()
	if err != nil {
		log.Fatal(err)
	}
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func validBlock(block, prevBlock *Block) bool {
	// Confirm the hashes
	if !bytes.Equal(block.PrevBlockHash, block.Hash) {
		return false
	}
	// confirm the block's hash is valid
	if !block.validateHash(block.Hash) {
		return false
	}

	/*
		if prevBlock.Pos+1 != block.Pos {
			return false
		}
	*/
	return true
}

func (b *Block) validateHash(hash []byte) bool {
	b.SetHash()
	if !bytes.Equal(b.Hash, hash) {
		return false
	}
	return true
}

func ViewAllFromStore() {
	db := logstore.Open()
	defer db.Close()

	_ = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("block"))

		_ := b.ForEach(func(k, v []byte) error {

			var b Block
			json.Unmarshal(v, &b)

			fmt.Printf("Pos: %x\n", b.Pos)
			fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
			fmt.Printf("Data: %s\n", b.Data)
			fmt.Printf("Hash: %x\n", b.Hash)
			fmt.Println()

			//spew.Dump(b)
			return nil
		})
		return nil
	})
}
