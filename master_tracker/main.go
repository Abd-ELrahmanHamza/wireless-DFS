package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jwangsadinata/go-multimap/slicemultimap"
)

type DataNode struct {
	ID            int32
	Addr          string // IP:Port
	LastPingstamp time.Time
}

// String function for DataNode
func (d *DataNode) String() string {
	return fmt.Sprintf("ID: %v, Addr: %v", d.ID, d.Addr)
}
func (d *DataNode) isAlive() bool {
	return time.Since(d.LastPingstamp) < time.Second
}

type lookupEntry struct {
	DataKeeperNode *DataNode
	filePath       string
	isAlive        bool
}

// implement print function for lookupTableEntry
func (l *lookupEntry) String() string {
	return fmt.Sprintf("DataKeeperNode: %v, filePath: %v, isAlive: %v", l.DataKeeperNode, l.filePath, l.DataKeeperNode.isAlive())
}

var (
	lookupTable   = slicemultimap.New()       // schema is file name -> lookupEntry
	DataNodes_Map = make(map[int32]*DataNode) // map of data keeper nodes with the data keeper node ID as the key
)

func main() {
	// create a lookup table of file names and their corresponding data keeper nodes with the file name as the key
	lookupTable.Put("file2", &lookupEntry{&DataNode{3, "localhost:8083", time.Now()}, "file2", true})
	log.Println(lookupTable.Get("file2"))

	println(nodesCounter())
	println(nodesCounter())
	// run_grpc()
}
