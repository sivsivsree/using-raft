package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/sivsivsree/using-raft/logstore"
	bolt "go.etcd.io/bbolt"
	"log"
	"strconv"
	"sync"
	"time"
)

var mux = sync.Mutex{}

type Block struct {
	Pos           int
	Timestamp     int64
	Data          string
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, []byte(b.Data), timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func NewBlock(data string, prevBlock *Block) *Block {
	pos := prevBlock.Pos + 1
	block := &Block{pos, time.Now().Unix(), data, prevBlock.Hash, []byte{}}
	block.SetHash()
	return block
}

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	// AddBlockToTheBucket

	db := logstore.Open()

	// key := strconv.Itoa(pos)
	er := db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("block"))

		prevBlock := bc.Blocks[len(bc.Blocks)-1]
		newBlock := NewBlock(data, prevBlock)

		newBlockBytes, _ := json.Marshal(newBlock)

		key := itob(newBlock.Pos)

		err := b.Put([]byte(key), newBlockBytes)

		bc.Blocks = append(bc.Blocks, newBlock)
		return err
	})

	db.Close()

	if er != nil {
		log.Fatal(er)
	}
}

func NewGenesisBlock() *Block {
	// Add NewGenesis to the Block

	db := logstore.Open()

	pos := 0
	genesis := &Block{pos, time.Now().Unix(), "Geneisis", []byte(""), []byte{}}
	genesis.SetHash()

	newBlockBytes, _ := json.Marshal(genesis)

	// key := strconv.Itoa(pos)
	er := db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("block"))

		fmt.Println("NewGenesisBlock")

		key := itob(genesis.Pos)
		err := b.Put([]byte(key), newBlockBytes)

		return err
	})

	db.Close()

	if er != nil {
		log.Fatal(er)
	}

	return genesis

}

func NewBlockchain() *Blockchain {
	// createBucketIfNotExists

	err := logstore.CreateBucketIfNotExist()
	if err != nil {
		log.Fatal(err)
	}

	// if data exists in genesis or not
	bc := &Blockchain{[]*Block{}}
	db := logstore.Open()

	err = db.View(func(tx *bolt.Tx) error {

		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("block"))

		err = b.ForEach(func(k, v []byte) error {

			var b Block
			_ = json.Unmarshal(v, &b)

			bc.Blocks = append(bc.Blocks, &b)

			return nil
		})
		return err
	})

	db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if len(bc.Blocks) == 0 {
		bc = &Blockchain{[]*Block{NewGenesisBlock()}}
	}

	return bc
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

	_ = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("block"))

		_ = b.ForEach(func(k, v []byte) error {

			var b Block
			_ = json.Unmarshal(v, &b)

			fmt.Printf("Pos: %d\n", b.Pos)
			fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
			fmt.Printf("Data: %s\n", b.Data)
			fmt.Printf("Hash: %x\n", b.Hash)
			fmt.Println()

			//spew.Dump(b)
			return nil
		})
		return nil
	})

	db.Close()
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
