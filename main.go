package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sivsivsree/using-raft/block"
	bolt "go.etcd.io/bbolt"
	"log"
)

func main() {

	db, err := bolt.Open("bolt.siv", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	chain, err := block.InitGenesis()
	_, _ = block.AddNewBlock("Hello there")
	_, _ = block.AddNewBlock("Hello therea")
	_, _ = block.AddNewBlock("Hello thereaa")
	_, _ = block.AddNewBlock("Hello thereaaa")
	_, _ = block.AddNewBlock("Hello thereaaaa")
	_, _ = block.AddNewBlock("Hello thereaaaaaa")
	_, _ = block.AddNewBlock("Hello thereaaaaaaaaed")

	defer db.Close()

	/*	err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		_ = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			v := b.Get([]byte("answer"))
			fmt.Printf("The answer is: %s\n", v)
			return nil
		})

		var a string
		prev := db.Stats()

		for {
			// Wait for 10s.
			_, _ = fmt.Scanf("%s", &a)

			// Grab the current stats and diff them.
			stats := db.Stats()
			diff := stats.Sub(&prev)

			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("MyBucket"))
				err := b.Put([]byte("answer"), []byte(a))
				return err
			})

			// Encode stats to JSON and print to STDERR.
			_ = json.NewEncoder(os.Stderr).Encode(diff)



			// Save stats for the next loop.
			prev = stats
		}*/

	spew.Dump(chain)

}
