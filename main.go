package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sivsivsree/using-raft/httpd"
	"github.com/sivsivsree/using-raft/store"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// Command line defaults
const (
	DefaultHTTPAddr = ":11000"
	DefaultRaftAddr = ":12000"
)

// Command line parameters
var inmem bool
var httpAddr string
var raftAddr string
var joinAddr string
var nodeID string

func init() {
	flag.BoolVar(&inmem, "inmem", false, "Use in-memory storage for Raft")
	flag.StringVar(&httpAddr, "haddr", DefaultHTTPAddr, "Set the HTTP bind address")
	flag.StringVar(&raftAddr, "raddr", DefaultRaftAddr, "Set Raft bind address")
	flag.StringVar(&joinAddr, "join", "", "Set join address, if any")
	flag.StringVar(&nodeID, "id", "", "Node ID")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <raft-data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}

	// Ensure Raft storage exists.
	raftDir := flag.Arg(0)
	if raftDir == "" {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}
	os.MkdirAll(raftDir, 0700)

	s := store.New(inmem)
	s.RaftDir = raftDir
	s.RaftBind = raftAddr
	if err := s.Open(joinAddr == "", nodeID); err != nil {
		log.Fatalf("failed to open store: %s", err.Error())
	}

	h := httpd.New(httpAddr, s)
	if err := h.Start(); err != nil {
		log.Fatalf("failed to start HTTP service: %s", err.Error())
	}

	// If join was specified, make the join request.
	if joinAddr != "" {
		if err := join(joinAddr, raftAddr, nodeID); err != nil {
			log.Fatalf("failed to join node at %s: %s", joinAddr, err.Error())
		}
	}

	log.Println("hraftd started successfully")

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("hraftd exiting")
}

func join(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr), "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

/*func main() {

port := flag.String("port", "7000", "Service port")

flag.Parse()

fmt.Println(*port)
bc := block.NewBlockchain()

//for i:= 0; i<100; i++ {
//	bc.AddBlock("data." + strconv.Itoa(i))
//}

//for _, b := range bc.Blocks {
//	fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
//	fmt.Printf("Data: %s\n", b.Data)
//	fmt.Printf("Hash: %x\n", b.Hash)
//	fmt.Println()
//}

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
	}

//spew.Dump(bc)

//block.ViewAllFromStore()

r := mux.NewRouter()

r.HandleFunc("/add/{data}", func (writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	data := vars["data"]
	bc.AddBlock(data)

	json.NewEncoder(writer).Encode(bc.Blocks)

}).Methods("GET")

r.HandleFunc("/get", func (writer http.ResponseWriter, request *http.Request) {
	block.ViewAllFromStore()

	json.NewEncoder(writer).Encode(bc.Blocks)

}).Methods("GET")

http.ListenAndServe(fmt.Sprintf(":%s", *port), r)
}*/
