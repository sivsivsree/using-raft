package logstore

import (
	bolt "go.etcd.io/bbolt"
	"log"
)

func Open() *bolt.DB {
	db, err := bolt.Open("genesis.siv", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func CreateBucketIfNotExist() error {
	db := Open()
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("block"))
		if err != nil {
			return err
		}
		return nil
	})
}
